package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Match struct {
	ID       string        `db:"id"`
	MatchID  string        `db:"match_id"`
	Version  string        `db:"version"`
	Date     time.Time     `db:"date"`
	Duration time.Duration `db:"duration"`
	WinnerID string        `db:"winner_id"`
}

type CreateMatch struct {
	MatchID  string
	Version  string
	Date     time.Time
	Duration time.Duration
	WinnerID string
}

type UpdateMatch struct {
	Version  string
	Date     time.Time
	Duration time.Duration
	WinnerID string
}

type MatchStore struct {
	tx DBTX
}

func NewMatchStore(tx DBTX) *MatchStore {
	return &MatchStore{
		tx: tx,
	}
}

func (s *MatchStore) GetMatchByMatchID(ctx context.Context, matchID string) (*Match, error) {
	rows, err := s.tx.Query(ctx, `
		select id, match_id, version, date, duration, winner_id
		from Match
		where match_id = @match_id;
	`, pgx.StrictNamedArgs{
		"match_id": matchID,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Match])
}

func (s *MatchStore) CreateMatch(ctx context.Context, req CreateMatch) (*Match, error) {
	rows, err := s.tx.Query(ctx, `
		insert into Match (match_id, version, date, duration, winner_id)
		values (@match_id, @version, @date, @duration, @winner_id)
		returning id, match_id, version, date, duration, winner_id;
	`, pgx.StrictNamedArgs{
		"match_id":  req.MatchID,
		"version":   req.Version,
		"date":      req.Date,
		"duration":  req.Duration,
		"winner_id": req.WinnerID,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Match])
}

func (s *MatchStore) UpdateMatchByMatchID(ctx context.Context, matchID string, req UpdateMatch) (*Match, error) {
	rows, err := s.tx.Query(ctx, `
		update Match
		set version = @version, date = @date, duration = @duration, winner_id = @winner_id
		where match_id = @match_id
		returning id, match_id, version, date, duration, winner_id;
	`, pgx.StrictNamedArgs{
		"match_id":  matchID,
		"version":   req.Version,
		"date":      req.Date,
		"duration":  req.Duration,
		"winner_id": req.WinnerID,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Match])
}

func (s *MatchStore) DeleteMatchByMatchID(ctx context.Context, matchID string) (*Match, error) {
	rows, err := s.tx.Query(ctx, `
		delete from Match
		where match_id = @match_id
		returning id, match_id, version, date, duration, winner_id;
	`, pgx.StrictNamedArgs{
		"match_id": matchID,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Match])
}
