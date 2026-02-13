package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *DB {
	return &DB{pool: pool}
}

func (db *DB) RecordProfile(ctx context.Context, summoner *internal.Profile) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var (
		rankStore     = postgres.RankStore{Tx: tx}
		summonerStore = postgres.SummonerStore{Tx: tx}
	)

	summonerIn := postgres.Summoner{
		PUUID:   summoner.PUUID,
		Name:    summoner.Name,
		Tagline: summoner.Tagline,
	}

	err = summonerStore.CreateSummoner(ctx, summonerIn)
	if err != nil {
		return err
	}

	rankStatus := postgres.RankStatus{
		PUUID:         summoner.PUUID.String(),
		EffectiveDate: summoner.Rank.EffectiveDate,
		IsRanked:      false,
	}

	if summoner.Rank.Detail != nil {
		rankStatus.IsRanked = true
	}

	statusID, err := rankStore.CreateRankStatus(ctx, rankStatus)
	if err != nil {
		return err
	}

	if summoner.Rank.Detail != nil {
		rankDetail := postgres.RankDetail{
			RankStatusID: statusID,
			Wins:         summoner.Rank.Detail.Wins,
			Losses:       summoner.Rank.Detail.Wins,
			Tier:         summoner.Rank.Detail.Rank.Tier.String(),
			Division:     summoner.Rank.Detail.Rank.Division.String(),
			LP:           summoner.Rank.Detail.Rank.LP,
		}

		if err := rankStore.CreateRankDetail(ctx, rankDetail); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (db *DB) GetProfile(ctx context.Context, puuid riot.PUUID) (m *internal.Profile, err error) {
	summonerStore := postgres.SummonerStore{Tx: db.pool}

	summoner, err := summonerStore.GetSummoner(ctx, puuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err := fmt.Errorf("%w: %w", internal.ErrSummonerNotFound, err)
			return m, err
		}

		return nil, err
	}

	status, detail, err := (*DB)(db).getMostRecentRank(ctx, puuid)
	if err != nil {
		return nil, err
	}

	mapper := postgresToRankStatus{
		Status: status,
		Detail: detail,
	}

	profile := internal.Profile{
		PUUID:   summoner.PUUID,
		Name:    summoner.Name,
		Tagline: summoner.Tagline,
		Rank:    mapper.Map(),
	}

	return &profile, nil
}

func (db *DB) SearchSummoner(ctx context.Context, q string) ([]internal.SearchResult, error) {
	var (
		summonerStore = postgres.SummonerStore{Tx: db.pool}
		rankStore     = postgres.RankStore{Tx: db.pool}
	)

	storeResults, err := summonerStore.SearchSummoner(ctx, q)
	if err != nil {
		return nil, err
	}

	mostRecentStatusIDs := []*int{}
	for _, result := range storeResults {
		ids, err := rankStore.ListRankIDs(ctx, result.PUUID, postgres.ListRankOption{Limit: 1, Recent: true})
		if err != nil {
			return nil, err
		}

		if len(ids) == 0 {
			mostRecentStatusIDs = append(mostRecentStatusIDs, nil)
		} else {
			mostRecentStatusIDs = append(mostRecentStatusIDs, &ids[0])
		}
	}

	statusList := []*postgres.RankStatus{}
	detailList := []*postgres.RankDetail{}
	for _, id := range mostRecentStatusIDs {
		if id == nil {
			statusList = append(statusList, nil)
			detailList = append(detailList, nil)
			continue
		}

		status, err := rankStore.GetRankStatus(ctx, *id)
		if err != nil {
			return nil, err
		}

		detail, err := rankStore.GetRankDetail(ctx, *id)
		if err != nil {
			return nil, err
		}

		statusList = append(statusList, &status)
		detailList = append(detailList, &detail)
	}

	results := []internal.SearchResult{}
	for i := range len(storeResults) {
		results = append(results, toSearchResult(storeResults[i], statusList[i], detailList[i]))
	}

	return results, nil
}

func (db *DB) SearchByNameTag(ctx context.Context, name, tag string) ([]internal.Profile, error) {
	summoner := postgres.SummonerStore{Tx: db.pool}

	result, err := summoner.SearchByNameTag(ctx, name, tag)
	if err != nil {
		return nil, err
	}

	profiles := []internal.Profile{}
	for _, summoner := range result {
		profiles = append(profiles, internal.Profile{
			PUUID:   summoner.PUUID,
			Name:    summoner.Name,
			Tagline: summoner.Tagline,
		})
	}

	return profiles, nil
}

func toSearchResult(summoner postgres.Summoner, status *postgres.RankStatus, detail *postgres.RankDetail) internal.SearchResult {
	result := internal.SearchResult{
		PUUID:   summoner.PUUID,
		Name:    summoner.Name,
		Tagline: summoner.Tagline,
		Rank:    nil,
	}

	if status != nil {
		result.Rank = &internal.RankStatus{
			PUUID:         summoner.PUUID,
			EffectiveDate: status.EffectiveDate,
			Detail:        nil,
		}

		if detail != nil {
			result.Rank.Detail = &internal.RankDetail{
				Wins:   detail.Wins,
				Losses: detail.Losses,
				Rank: internal.Rank{
					Tier:     riot.Tier(detail.Tier),
					Division: riot.Division(detail.Division),
					LP:       detail.LP,
				},
			}
		}
	}

	return result
}

func (db *DB) getMostRecentRank(ctx context.Context, puuid riot.PUUID) (m postgres.RankStatus, n *postgres.RankDetail, err error) {
	rankStore := postgres.RankStore{Tx: db.pool}

	ids, err := rankStore.ListRankIDs(ctx, puuid, postgres.ListRankOption{Limit: 1, Recent: true})
	if err != nil {
		return m, n, err
	}

	if len(ids) != 1 {
		return m, n, errors.New("ListRankIDS did not return exactly one id")
	}

	id := ids[0]

	return db.getRank(ctx, id)
}

func (db *DB) getRank(ctx context.Context, statusID int) (postgres.RankStatus, *postgres.RankDetail, error) {
	rankStore := postgres.RankStore{Tx: db.pool}

	status, err := rankStore.GetRankStatus(ctx, statusID)
	if err != nil {
		return postgres.RankStatus{}, nil, err
	}

	detail, err := rankStore.GetRankDetail(ctx, statusID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return status, nil, nil
		}

		return postgres.RankStatus{}, nil, err
	}

	return status, &detail, err
}

type postgresToRankStatus struct {
	Status postgres.RankStatus
	Detail *postgres.RankDetail
}

func (mapper postgresToRankStatus) Map() internal.RankStatus {
	status := mapper.Status
	detail := mapper.Detail

	result := internal.RankStatus{
		PUUID:         riot.PUUID(status.PUUID),
		EffectiveDate: status.EffectiveDate,
		Detail:        nil,
	}

	if detail != nil {
		result.Detail = &internal.RankDetail{
			Wins:   detail.Wins,
			Losses: detail.Losses,
			Rank: internal.Rank{
				Tier:     convertStringToRiotTier(detail.Tier),
				Division: convertStringToRiotRank(detail.Division),
				LP:       detail.LP,
			},
		}
	}

	return result
}

func convertStringToRiotTier(tier string) riot.Tier {
	switch tier {
	default:
		panic("bro.")
	case "Iron":
		return riot.TierIron
	case "Bronze":
		return riot.TierBronze
	case "Silver":
		return riot.TierSilver
	case "Gold":
		return riot.TierGold
	case "Platinum":
		return riot.TierPlatinum
	case "Emerald":
		return riot.TierEmerald
	case "Diamond":
		return riot.TierDiamond
	case "Master":
		return riot.TierMaster
	case "Grandmaster":
		return riot.TierGrandmaster
	case "Challenger":
		return riot.TierChallenger
	}
}

func convertStringToRiotRank(rank string) riot.Division {
	switch rank {
	default:
		panic("bro.")
	case "I":
		return riot.Division1
	case "II":
		return riot.Division2
	case "III":
		return riot.Division3
	case "IV":
		return riot.Division4
	}
}

func convertListToRunePage(runes [11]int) internal.RunePage {
	page := internal.RunePage{
		PrimaryTree:     runes[0],
		PrimaryKeystone: runes[1],
		PrimaryA:        runes[2],
		PrimaryB:        runes[3],
		PrimaryC:        runes[4],
		SecondaryTree:   runes[5],
		SecondaryA:      runes[6],
		SecondaryB:      runes[7],
		MiniOffense:     runes[8],
		MiniFlex:        runes[9],
		MiniDefense:     runes[10],
	}

	return page
}

func convertRunePageToList(runes internal.RunePage) [11]int {
	ids := [11]int{
		runes.PrimaryTree,
		runes.PrimaryKeystone,
		runes.PrimaryA,
		runes.PrimaryB,
		runes.PrimaryC,
		runes.SecondaryTree,
		runes.SecondaryA,
		runes.SecondaryB,
		runes.MiniOffense,
		runes.MiniFlex,
		runes.MiniDefense,
	}

	return ids
}

func convertTeamPositionToString(position internal.TeamPosition) string {
	switch position {
	case internal.TeamPositionBottom:
		return "Bottom"
	case internal.TeamPositionSupport:
		return "Support"
	case internal.TeamPositionTop:
		return "Top"
	case internal.TeamPositionMiddle:
		return "Middle"
	case internal.TeamPositionJungle:
		return "Jungle"
	}

	return ""
}

func convertStringToTeamPosition(s string) internal.TeamPosition {
	switch s {
	case "Bottom":
		return internal.TeamPositionBottom
	case "Support":
		return internal.TeamPositionSupport
	case "Middle":
		return internal.TeamPositionMiddle
	case "Top":
		return internal.TeamPositionTop
	case "Jungle":
		return internal.TeamPositionJungle
	}

	return internal.TeamPositionTop
}
