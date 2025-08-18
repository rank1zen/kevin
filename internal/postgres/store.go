package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

// Should be unexported
func CreateListRankOption(history []Match) [][2]ListRankOption {
	opt := [][2]ListRankOption{}

	for i, m := range history {
		rankBefore := ListRankOption{
			Offset: 0,
			Limit:  100,
			Recent: true,
		}

		rankAfter := ListRankOption{
			Offset: 0,
			Limit:  100,
			Recent: true,
		}

		if i > 0 {
			rankBefore.Start = &history[i-1].Date
		}

		rankBefore.End = &m.Date

		if i < len(history)-1 {
			rankAfter.End = &history[i+1].Date
		}

		rankAfter.Start = &m.Date

		opt = append(opt, [2]ListRankOption{rankBefore, rankAfter})
	}

	return opt
}

// Store manages connections with a postgres database.
type Store struct {
	Pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) internal.Store {
	return &Store{Pool: pool}
}

func (db *Store) GetChampions(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]internal.SummonerChampion, error) {
	panic("not implemented")
}

func (db *Store) SearchSummoner(ctx context.Context, q string) ([]internal.SearchResult, error) {
	summonerStore := SummonerStore{Tx: db.Pool}

	rankStore := RankStore{Tx: db.Pool}

	results, err := summonerStore.SearchSummoner(ctx, q)
	if err != nil {
		return nil, err
	}

	mostRecentStatusIDs := []*int{}
	for _, result := range results {
		ids, err := rankStore.ListRankIDs(ctx, result.PUUID, ListRankOption{Limit: 1, Recent: true})
		if err != nil {
			return nil, err
		}

		if len(ids) == 0 {
			mostRecentStatusIDs = append(mostRecentStatusIDs, nil)
		} else {
			mostRecentStatusIDs = append(mostRecentStatusIDs, &ids[0])
		}
	}

	statusList := []*RankStatus{}
	detailList := []*RankDetail{}
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

	dada := []internal.SearchResult{}
	for i := range len(results) {
		s := PostgresSearchResult2{
			Summoner: results[i],
			Status:   statusList[i],
			Detail:   detailList[i],
		}

		dada = append(dada, s.Convert())
	}

	return dada, nil
}

func (db *Store) RecordProfile(ctx context.Context, summoner internal.Profile) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	rankStore := RankStore{Tx: tx}

	summonerStore := SummonerStore{Tx: tx}

	summonerIn := Summoner{
		PUUID:   summoner.PUUID,
		Name:    summoner.Name,
		Tagline: summoner.Tagline,
	}

	err = summonerStore.CreateSummoner(ctx, summonerIn)
	if err != nil {
		return err
	}

	rankStatus := RankStatus{
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
		rankDetail := RankDetail{
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

	tx.Commit(ctx)

	return nil
}

func (db *Store) GetProfileDetail(ctx context.Context, puuid riot.PUUID) (internal.ProfileDetail, error) {
	summonerStore := SummonerStore{Tx: db.Pool}

	summoner, err := summonerStore.GetSummoner(ctx, puuid)
	if err != nil {
		return internal.ProfileDetail{}, err
	}

	detail := internal.ProfileDetail{
		PUUID:   summoner.PUUID,
		Name:    summoner.Name,
		Tagline: summoner.Tagline,
		Rank:    internal.RankStatus{},
	}

	return detail, nil
}

func (db *Store) RecordMatch(ctx context.Context, match internal.Match) error {
	matchStore := MatchStore{Tx: db.Pool}

	batch := pgx.Batch{}

	matchStore.CreateMatchInBatch(&batch, Match{
		ID:       match.ID,
		Date:     match.Date,
		Duration: match.Duration,
		Version:  match.Version,
		WinnerID: match.WinnerID,
	})

	for _, p := range match.Participants {
		matchStore.CreateParticipantInBatch(&batch, Participant{
			PUUID:                p.PUUID.String(),
			MatchID:              p.MatchID,
			TeamID:               p.TeamID,
			ChampionID:           p.ChampionID,
			ChampionLevel:        p.ChampionLevel,
			TeamPosition:         convertTeamPositionToString(p.TeamPosition),
			SummonerIDs:          p.SummonerIDs,
			Runes:                convertRunePageToList(p.Runes),
			Items:                p.Items,
			Kills:                p.Kills,
			Deaths:               p.Deaths,
			Assists:              p.Assists,
			KillParticipation:    p.KillParticipation,
			CreepScore:           p.CreepScore,
			CreepScorePerMinute:  p.CreepScorePerMinute,
			DamageDealt:          p.DamageDealt,
			DamageTaken:          p.DamageTaken,
			DamageDeltaEnemy:     p.DamageDeltaEnemy,
			DamagePercentageTeam: p.DamagePercentageTeam,
			GoldEarned:           p.GoldEarned,
			GoldDeltaEnemy:       p.GoldDeltaEnemy,
			GoldPercentageTeam:   p.GoldPercentageTeam,
			VisionScore:          p.VisionScore,
			PinkWardsBought:      p.PinkWardsBought,
		})
	}

	br := db.Pool.SendBatch(ctx, &batch)

	return br.Close()
}

func (db *Store) GetMatchDetail(ctx context.Context, id string) (internal.MatchDetail, error) {
	m := internal.MatchDetail{}

	matchStore := MatchStore{Tx: db.Pool}

	summonerStore := SummonerStore{Tx: db.Pool}

	match, err := matchStore.GetMatch(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return m, fmt.Errorf("%w: %w", internal.ErrMatchNotFound, err)
		}

		return m, fmt.Errorf("%w: %w", internal.ErrUnknownStoreError, err)
	}

	participants, err := matchStore.GetParticipants(ctx, id)
	if err != nil {
		return m, fmt.Errorf("%w: %w", internal.ErrUnknownStoreError, err)
	}

	if len(participants) != 10 {
		return m, internal.ErrMatchMissingParticipants
	}

	detail := internal.MatchDetail{
		ID:           id,
		Date:         match.Date,
		Duration:     match.Duration,
		Version:      match.Version,
		WinnerID:     match.WinnerID,
		Participants: [10]internal.ParticipantDetail{},
	}

	for i := range 10 {
		puuid := participants[i].PUUID
		summoner, err := summonerStore.GetSummoner(ctx, internal.NewPUUIDFromString(puuid))
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return m, err
			}

			detail.Participants[i] = internal.NewParticipantDetail(ParticipantDetailFromPG(participants[i], Summoner{PUUID: internal.NewPUUIDFromString(puuid)}, nil, nil, nil))
			continue
		}

		detail.Participants[i] = internal.NewParticipantDetail(ParticipantDetailFromPG(participants[i], summoner, nil, nil, nil))
	}

	return detail, nil
}

func (db *Store) GetMatchHistory(ctx context.Context, puuid riot.PUUID, start, end time.Time) (matches []internal.SummonerMatch, err error) {
	defer errWrap(&err, "GetMatchHistory")

	matchStore := MatchStore{Tx: db.Pool}

	rankStore := RankStore{Tx: db.Pool}

	ids, err := matchStore.ListMatchHistoryIDs(ctx, puuid, start, end)
	if err != nil {
		return nil, err
	}

	matchHistory := []internal.SummonerMatch{}

	var (
		matchList       []Match       = []Match{}
		participantList []Participant = []Participant{}
		rankBeforeList  []*RankFull   = []*RankFull{}
		rankAfterList   []*RankFull   = []*RankFull{}
	)

	for _, id := range ids {
		match, err := matchStore.GetMatch(ctx, id)
		if err != nil {
			return nil, err
		}

		matchList = append(matchList, match)

		participant, err := matchStore.GetParticipant(ctx, puuid, id)
		if err != nil {
			return nil, err
		}

		participantList = append(participantList, participant)
	}

	listRankOptions := CreateListRankOption(matchList)
	for _, opt := range listRankOptions {
		statusIDs, err := rankStore.ListRankIDs(ctx, puuid, opt[0])
		if err != nil {
			return nil, err
		}

		id := chooseStatusID(statusIDs)
		if id != nil {
			status, detail, err := db.getRank(ctx, *id)
			if err != nil {
				return nil, err
			}

			rankBeforeList = append(rankBeforeList, &RankFull{Status: status, Detail: detail})
		} else {
			rankBeforeList = append(rankBeforeList, nil)
		}

		statusIDs, err = rankStore.ListRankIDs(ctx, puuid, opt[1])
		if err != nil {
			return nil, err
		}

		id = chooseStatusID(statusIDs)
		if id != nil {
			status, detail, err := db.getRank(ctx, *id)
			if err != nil {
				return nil, err
			}

			rankAfterList = append(rankAfterList, &RankFull{Status: status, Detail: detail})
		} else {
			rankAfterList = append(rankAfterList, nil)
		}
	}

	for i := range ids {
		match := internal.NewSummonerMatch(WithPostgresSummonerMatch(matchList[i], participantList[i], rankBeforeList[i], rankAfterList[i]))
		matchHistory = append(matchHistory, match)
	}

	return matchHistory, nil
}

func (db *Store) GetNewMatchIDs(ctx context.Context, ids []string) (newIDs []string, err error) {
	panic("not implemented")
}

func ParticipantDetailFromPG(participant Participant, summoner Summoner, currentRank, rankBefore, rankAfter *RankFull) internal.ParticipantDetailOption {
	return func(m *internal.ParticipantDetail) {
		m.PUUID = internal.NewPUUIDFromString(participant.PUUID)

		m.MatchID = participant.MatchID

		m.TeamID = participant.TeamID

		m.ChampionID = participant.ChampionID

		m.ChampionLevel = participant.ChampionLevel

		m.TeamPosition = convertStringToTeamPosition(participant.TeamPosition)

		m.SummonerIDs = participant.SummonerIDs

		m.Runes = convertListToRunePage(participant.Runes)

		m.Items = participant.Items

		m.Kills = participant.Kills

		m.Deaths = participant.Deaths

		m.Assists = participant.Assists

		m.KillParticipation = participant.KillParticipation

		m.CreepScore = participant.CreepScore

		m.CreepScorePerMinute = participant.CreepScorePerMinute

		m.DamageDealt = participant.DamageDealt

		m.DamageTaken = participant.DamageTaken

		m.DamageDeltaEnemy = participant.DamageDeltaEnemy

		m.DamagePercentageTeam = participant.DamagePercentageTeam

		m.GoldEarned = participant.GoldEarned

		m.GoldDeltaEnemy = participant.GoldDeltaEnemy

		m.GoldPercentageTeam = participant.GoldPercentageTeam

		m.VisionScore = participant.VisionScore

		m.PinkWardsBought = participant.PinkWardsBought

		m.Name = summoner.Name

		m.Tag = summoner.Tagline

		m.RankAfter = nil
		if rankAfter != nil {
			rank := internal.NewRankStatus(WithPostgresRankStatus2(rankAfter))
			m.RankBefore = &rank
		}

		m.RankBefore = nil
		if rankBefore != nil {
			rank := internal.NewRankStatus(WithPostgresRankStatus2(rankBefore))
			m.RankBefore = &rank
		}

		m.CurrentRank = nil
		if currentRank != nil {
			rank := internal.NewRankStatus(WithPostgresRankStatus2(currentRank))
			m.RankBefore = &rank
		}
	}
}

// Deprecated: not using.
func WithPostgresSummonerMatch(match Match, participant Participant, rankBefore, rankAfter *RankFull) internal.SummonerMatchOption {
	return func(m *internal.SummonerMatch) {
		m.PUUID = internal.NewPUUIDFromString(participant.PUUID)

		m.MatchID = match.ID

		m.TeamID = participant.TeamID

		m.ChampionID = participant.ChampionID

		m.ChampionLevel = participant.ChampionLevel

		m.TeamPosition = convertStringToTeamPosition(participant.TeamPosition)

		m.SummonerIDs = participant.SummonerIDs

		m.Runes = convertListToRunePage(participant.Runes)

		m.Items = participant.Items

		m.Kills = participant.Kills

		m.Deaths = participant.Deaths

		m.Assists = participant.Assists

		m.KillParticipation = participant.KillParticipation

		m.CreepScore = participant.CreepScore

		m.CreepScorePerMinute = participant.CreepScorePerMinute

		m.DamageDealt = participant.DamageDealt

		m.DamageTaken = participant.DamageTaken

		m.DamageDeltaEnemy = participant.DamageDeltaEnemy

		m.DamagePercentageTeam = participant.DamagePercentageTeam

		m.GoldEarned = participant.GoldEarned

		m.GoldDeltaEnemy = participant.GoldDeltaEnemy

		m.GoldPercentageTeam = participant.GoldPercentageTeam

		m.VisionScore = participant.VisionScore

		m.PinkWardsBought = participant.PinkWardsBought

		m.Date = match.Date

		m.Duration = match.Duration

		m.Win = false
		if participant.TeamID == match.WinnerID {
			m.Win = true
		}

		m.RankBefore = nil
		if rankBefore != nil {
			rank := internal.NewRankStatus(WithPostgresRankStatus2(rankBefore))
			m.RankBefore = &rank
		}

		m.RankAfter = nil
		if rankAfter != nil {
			rank := internal.NewRankStatus(WithPostgresRankStatus2(rankAfter))
			m.RankAfter = &rank
		}
	}
}

// Deprecated: not using.
func WithPostgresRankStatus2(rank *RankFull) internal.RankStatusOption {
	return func(m *internal.RankStatus) {
		m.PUUID = internal.NewPUUIDFromString(rank.Status.PUUID)
		m.EffectiveDate = rank.Status.EffectiveDate

		m.Detail = nil
		if detail := rank.Detail; detail != nil {
			m.Detail = &internal.RankDetail{
				Wins:   detail.Wins,
				Losses: detail.Losses,
				Rank: internal.Rank{
					Tier:     convertStringToRiotTier(detail.Tier),
					Division: convertStringToRiotRank(detail.Division),
					LP:       detail.LP,
				},
			}
		}
	}
}

type PostgresSearchResult2 struct {
	Summoner Summoner

	Status *RankStatus
	Detail *RankDetail
}

func (m *PostgresSearchResult2) Convert() internal.SearchResult {
	result := internal.SearchResult{
		PUUID:   m.Summoner.PUUID,
		Name:    m.Summoner.Name,
		Tagline: m.Summoner.Tagline,
		Rank:    nil,
	}

	if m.Status != nil {
		result.Rank = &internal.RankStatus{
			PUUID:         m.Summoner.PUUID,
			EffectiveDate: m.Status.EffectiveDate,
			Detail:        nil,
		}

		if m.Detail != nil {
			result.Rank.Detail = &internal.RankDetail{
				Wins:   m.Detail.Wins,
				Losses: m.Detail.Losses,
				Rank: internal.Rank{
					Tier:     riot.Tier(m.Detail.Tier),
					Division: riot.Division(m.Detail.Division),
					LP:       m.Detail.LP,
				},
			}
		}
	}

	return result
}

// chooseStatusID chooses some id that is suitable.
func chooseStatusID(statusIDs []int) *int {
	if len(statusIDs) == 0 {
		return nil
	}
	return &statusIDs[0]
}

func (db *Store) getRank(ctx context.Context, statusID int) (RankStatus, *RankDetail, error) {
	rankStore := RankStore{Tx: db.Pool}

	status, err := rankStore.GetRankStatus(ctx, statusID)
	if err != nil {
		return RankStatus{}, nil, err
	}

	detail, err := rankStore.GetRankDetail(ctx, statusID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return status, nil, nil
		}

		return RankStatus{}, nil, err
	}

	return status, &detail, err
}
