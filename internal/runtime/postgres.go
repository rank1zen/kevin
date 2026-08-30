package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initializePostgres(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	pool, err := retry.NewWithData[*pgxpool.Pool](
		retry.Attempts(5),
		retry.Delay(100*time.Millisecond),
	).Do(func() (*pgxpool.Pool, error) {
		pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create postgres connection pool: %w", err)
		}

		if err := pool.Ping(ctx); err != nil {
			return nil, fmt.Errorf("failed to ping postgres: %w", err)
		}

		return pool, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize postgres: %w", err)
	}
	return pool, nil
}

func initializePgxConn(ctx context.Context, cfg Config) (*pgx.Conn, error) {
	pgxCfg, err := pgx.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	conn, err := retry.NewWithData[*pgx.Conn](
		retry.Attempts(5),
		retry.Delay(100*time.Millisecond),
	).Do(func() (*pgx.Conn, error) {
		conn, err := pgx.ConnectConfig(ctx, pgxCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create postgres connection: %w", err)
		}

		if err := conn.Ping(ctx); err != nil {
			return nil, fmt.Errorf("failed to ping postgres: %w", err)
		}

		return conn, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize postgres: %w", err)
	}

	return conn, nil
}
