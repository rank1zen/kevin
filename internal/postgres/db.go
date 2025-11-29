package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, connString string) (internal.DB, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres connection pool: %w", err)
	}

	return &DB{pool: pool}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}
