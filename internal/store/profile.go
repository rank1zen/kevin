package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
)

type ProfileStore Store

func (db *ProfileStore) RecordProfile(ctx context.Context, summoner *internal.Profile) error {
	tx, err := db.Pool.Begin(ctx)
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

func (db *ProfileStore) GetProfile(ctx context.Context, puuid riot.PUUID) (m *internal.Profile, err error) {
	summonerStore := postgres.SummonerStore{Tx: db.Pool}

	summoner, err := summonerStore.GetSummoner(ctx, puuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err := fmt.Errorf("%w: %w", internal.ErrSummonerNotFound, err)
			return m, err
		}

		return nil, err
	}

	status, detail, err := (*Store)(db).getMostRecentRank(ctx, puuid)
	if err != nil {
		return nil, err
	}

	mapper := PostgresToRankStatus{
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

func (db *ProfileStore) SearchSummoner(ctx context.Context, q string) ([]internal.SearchResult, error) {
	var (
		summonerStore = postgres.SummonerStore{Tx: db.Pool}
		rankStore     = postgres.RankStore{Tx: db.Pool}
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

func (db *ProfileStore) SearchByNameTag(ctx context.Context, name, tag string) ([]internal.Profile, error) {
	summoner := postgres.SummonerStore{Tx: db.Pool}

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
