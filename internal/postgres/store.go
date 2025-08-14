package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

// TODO: replace Store.
type Store2 struct {
	Pool *pgxpool.Pool
}

func NewStore2() internal.Store2 {
	return &Store2{}
}

func (db *Store2) RecordProfile(ctx context.Context, summoner internal.Profile) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	rank := RankStore{Tx: tx}

	// create or update summoner

	rankStatus := RankStatus{
		PUUID:         summoner.PUUID.String(),
		EffectiveDate: summoner.Rank.EffectiveDate,
		IsRanked:      false,
	}

	if summoner.Rank.Detail != nil {
		rankStatus.IsRanked = true
	}

	statusID, err := rank.CreateRankStatus(ctx, rankStatus)
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

		if err := rank.CreateRankDetail(ctx, rankDetail); err != nil {
			return err
		}
	}

	tx.Commit(ctx)

	return nil
}

func (db *Store2) GetProfileDetail(ctx context.Context, puuid riot.PUUID) (internal.ProfileDetail, error) {
	panic("not implemented")
}

func (db *Store2) RecordMatch(ctx context.Context, match internal.Match) error {
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

func (db *Store2) GetMatchDetail(ctx context.Context, id riot.PUUID) (internal.MatchDetail, error) {
	panic("not implemented")
}

func (db *Store2) GetMatchHistory(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]internal.SummonerMatch, error) {
	// list games played for puuid
	// get participant for each game
	// get match for each game
	// list ranks for puuid between games
	// get rank record for that id

	panic("not implemented")
}

func (db *Store2) GetNewMatchIDs(ctx context.Context, ids []string) (newIDs []string, err error) {
	panic("not implemented")
}
