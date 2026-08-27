package store

import (
	"context"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
)

type Rank struct {
	ID           uuid.UUID `db:"id"`
	PUUID        string    `db:"puuid"`
	Wins         int       `db:"wins"`
	Losses       int       `db:"losses"`
	Tier         string    `db:"tier"`
	Division     string    `db:"division"`
	LeaguePoints int       `db:"league_points"`
	LastUpdated  time.Time `db:"last_updated"`
}

type CreateRank struct {
	PUUID        string
	Wins         int
	Losses       int
	Tier         string
	Division     string
	LeaguePoints int
	LastUpdated  time.Time
}

type UpdateRank struct {
	Wins         int
	Losses       int
	Tier         string
	Division     string
	LeaguePoints int
	LastUpdated  time.Time
}

type RankStore struct {
	tx DBTX
}

func NewRankStore(tx DBTX) *RankStore {
	return &RankStore{
		tx: tx,
	}
}

func (s *RankStore) GetRankByPUUID(ctx context.Context, puuid string) (*Rank, error) {
	rows, err := s.tx.Query(ctx, `
		select id, puuid, wins, losses, tier, division, league_points, last_updated
		from Rank
		where puuid = @puuid;
	`, pgx.StrictNamedArgs{
		"puuid": puuid,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Rank])
}

func (s *RankStore) CreateRank(ctx context.Context, req CreateRank) (*Rank, error) {
	rows, err := s.tx.Query(ctx, `
		insert into Rank (puuid, wins, losses, tier, division, league_points, last_updated)
		values (@puuid, @wins, @losses, @tier, @division, @league_points, @last_updated)
		returning id, puuid, wins, losses, tier, division, league_points, last_updated;
	`, pgx.StrictNamedArgs{
		"puuid":         req.PUUID,
		"wins":          req.Wins,
		"losses":        req.Losses,
		"tier":          req.Tier,
		"division":      req.Division,
		"league_points": req.LeaguePoints,
		"last_updated":  req.LastUpdated,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Rank])
}

func (s *RankStore) UpdateRankByPUUID(ctx context.Context, puuid string, req UpdateRank) (*Rank, error) {
	rows, err := s.tx.Query(ctx, `
		update Rank
		set wins = @wins, losses = @losses, tier = @tier, division = @division, league_points = @league_points, last_updated = @last_updated
		where puuid = @puuid
		returning id, puuid, wins, losses, tier, division, league_points, last_updated;
	`, pgx.StrictNamedArgs{
		"puuid":         puuid,
		"wins":          req.Wins,
		"losses":        req.Losses,
		"tier":          req.Tier,
		"division":      req.Division,
		"league_points": req.LeaguePoints,
		"last_updated":  req.LastUpdated,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Rank])
}

func (s *RankStore) GetRankByPUUIDs(ctx context.Context, puuids []string) ([]*Rank, error) {
	rows, err := s.tx.Query(ctx, `
		select id, puuid, wins, losses, tier, division, league_points, last_updated
		from Rank
		where puuid = ANY(@puuids);
	`, pgx.StrictNamedArgs{
		"puuids": puuids,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Rank])
}
