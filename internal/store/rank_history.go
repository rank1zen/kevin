package store

import (
	"context"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
)

type RankHistory struct {
	ID           uuid.UUID  `db:"id"`
	PUUID        string     `db:"puuid"`
	ValidFrom    time.Time  `db:"valid_from"`
	ValidTo      *time.Time `db:"valid_to"`
	Wins         int        `db:"wins"`
	Losses       int        `db:"losses"`
	Tier         string     `db:"tier"`
	Division     string     `db:"division"`
	LeaguePoints int        `db:"league_points"`
}

type CreateRankHistory struct {
	PUUID        string
	ValidFrom    time.Time
	ValidTo      *time.Time
	Wins         int
	Losses       int
	Tier         string
	Division     string
	LeaguePoints int
}

type UpdateRankHistory struct {
	ValidFrom    time.Time
	ValidTo      *time.Time
	Wins         int
	Losses       int
	Tier         string
	Division     string
	LeaguePoints int
}

type RankHistoryStore struct {
	tx DBTX
}

func NewRankHistoryStore(tx DBTX) *RankHistoryStore {
	return &RankHistoryStore{
		tx: tx,
	}
}

func (s *RankHistoryStore) GetRankHistoryByID(ctx context.Context, id uuid.UUID) (*RankHistory, error) {
	rows, err := s.tx.Query(ctx, `
		select id, puuid, valid_from, valid_to, wins, losses, tier, division, league_points
		from RankHistory
		where id = @id;
	`, pgx.StrictNamedArgs{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[RankHistory])
}

func (s *RankHistoryStore) CreateRankHistory(ctx context.Context, req CreateRankHistory) (*RankHistory, error) {
	rows, err := s.tx.Query(ctx, `
		insert into RankHistory (puuid, valid_from, valid_to, wins, losses, tier, division, league_points)
		values (@puuid, @valid_from, @valid_to, @wins, @losses, @tier, @division, @league_points)
		returning id, puuid, valid_from, valid_to, wins, losses, tier, division, league_points;
	`, pgx.StrictNamedArgs{
		"puuid":         req.PUUID,
		"valid_from":    req.ValidFrom,
		"valid_to":      req.ValidTo,
		"wins":          req.Wins,
		"losses":        req.Losses,
		"tier":          req.Tier,
		"division":      req.Division,
		"league_points": req.LeaguePoints,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[RankHistory])
}

func (s *RankHistoryStore) UpdateRankHistoryByID(ctx context.Context, id uuid.UUID, req UpdateRankHistory) (*RankHistory, error) {
	rows, err := s.tx.Query(ctx, `
		update RankHistory
		set valid_from = @valid_from, valid_to = @valid_to, wins = @wins, losses = @losses, tier = @tier, division = @division, league_points = @league_points
		where id = @id
		returning id, puuid, valid_from, valid_to, wins, losses, tier, division, league_points;
	`, pgx.StrictNamedArgs{
		"id":            id,
		"valid_from":    req.ValidFrom,
		"valid_to":      req.ValidTo,
		"wins":          req.Wins,
		"losses":        req.Losses,
		"tier":          req.Tier,
		"division":      req.Division,
		"league_points": req.LeaguePoints,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[RankHistory])
}

func (s *RankHistoryStore) DeleteRankHistoryByID(ctx context.Context, id uuid.UUID) (*RankHistory, error) {
	rows, err := s.tx.Query(ctx, `
		delete from RankHistory
		where id = @id
		returning id, puuid, valid_from, valid_to, wins, losses, tier, division, league_points;
	`, pgx.StrictNamedArgs{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[RankHistory])
}
