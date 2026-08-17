package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal/feature/profile"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *DB {
	return &DB{pool: pool}
}

// GetChampions implements profile.ProfileStore.
func (d *DB) GetChampions(ctx context.Context, puuid string, start time.Time, end time.Time) ([]profile.ChampionAverage, error) {
	panic("unimplemented")
}

// GetMatchlist implements profile.ProfileStore.
func (d *DB) GetMatchlist(ctx context.Context, puuid string, start time.Time, end time.Time) ([]profile.Match, error) {
	panic("unimplemented")
}

// GetProfile implements profile.ProfileStore.
func (d *DB) GetProfile(ctx context.Context, puuid string) (*profile.Profile, error) {
	panic("unimplemented")
}

// GetRankHistory implements profile.ProfileStore.
func (d *DB) GetRankHistory(ctx context.Context, puuid string, start time.Time, end time.Time) ([]profile.RankStatus, error) {
	panic("unimplemented")
}

// RecordProfile implements profile.ProfileStore.
func (d *DB) RecordProfile(ctx context.Context, profile *profile.Profile) error {
	panic("unimplemented")
}
