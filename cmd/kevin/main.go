package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal/config"
	"github.com/rank1zen/kevin/internal/frontend/server"
	"github.com/rank1zen/kevin/internal/log"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/service"
	"github.com/rank1zen/kevin/internal/store"
)

type Runtime struct {
	Logger       *slog.Logger
	RiotClient   *riot.Client
	PostgresConn *pgxpool.Pool
}

func main() {
	ctx := context.Background()

	rt := Runtime{}

	rt.Logger = log.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		rt.Logger.Error("failed to load configuration", "err", err)
		os.Exit(1)
	}

	pool, err := connectPostgres(ctx, cfg.GetPostgresConnection())
	if err != nil {
		rt.Logger.Error("error starting server", "err", err)
		os.Exit(1)
	}
	rt.PostgresConn = pool

	defer pool.Close()

	st := store.NewStore(pool)
	client := riot.NewClient(cfg.GetRiotAPIKey())
	datasource := service.NewService(client, st, pool)

	srvr := server.New(datasource, cfg.GetPort(), server.WithLogger(rt.Logger))

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		err := srvr.Open()
		slog.Default().Error("error starting server", "err", err)
	}()

	rt.Logger.Info("server started", "address", cfg.GetPort(), "environment", "prod")

	<-ctx.Done()

	rt.Logger.Info("shutting down server")

	os.Exit(0)
}

func connectPostgres(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres connection pool: %w", err)
	}

	return pool, nil
}
