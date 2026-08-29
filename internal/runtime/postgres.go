package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go/v5"
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
