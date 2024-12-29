package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/yujin/internal"
	"go.opentelemetry.io/otel/trace"
)

type DB struct {
	pool   *pgxpool.Pool
	tracer trace.Tracer
}

func NewDB(ctx context.Context, url string) (internal.Repository, error) {
	pgxCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres connection string: %w", err)
	}

	pgxCfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		return conn.Ping(ctx) == nil
	}

	pgxCfg.ConnConfig.Tracer = &tracer{}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	// FIXME: DB currently does not implement internal.Repository
	return &DB{
		pool: pool,
	}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}
