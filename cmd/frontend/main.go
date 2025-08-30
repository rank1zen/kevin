package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
)

func main() {
	config := Config{
		RiotAPIKey:         os.Getenv("KEVIN_RIOT_API_KEY"),
		PostgresConnection: os.Getenv("KEVIN_POSTGRES_CONNECTION"),
	}

	flag.BoolVar(&config.DevelopmentMode, "development-mode", false, "set mode to development")

	flag.StringVar(&config.Address, "address", "0.0.0.0:4001", "set server address")

	flag.Parse()

	ctx := context.Background()

	if err := config.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	os.Exit(0)
}

type Config struct {
	// PostgresConnection is used to connect to postgres. Format is
	// specified by [pgxpool.ParseConfig].
	PostgresConnection string

	// RiotAPIKey is required to access riot.
	RiotAPIKey string

	Address string

	DevelopmentMode bool
}

func (c *Config) Run(ctx context.Context) error {
	if c.RiotAPIKey == "" {
		return errors.New("riot api key not provided")
	}

	if c.PostgresConnection == "" {
		return errors.New("postgres connection not provided")
	}

	address := "0.0.0.0:4001"
	if c.Address != "" {
		address = c.Address
	}

	var logHandlerOptions slog.HandlerOptions
	if c.DevelopmentMode {
		logHandlerOptions.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &logHandlerOptions)).With("dev", c.DevelopmentMode)

	pool, err := connectPostgres(ctx, c.PostgresConnection)
	if err != nil {
		return err
	}

	defer pool.Close()

	store := postgres.NewStore(pool)

	client := riot.NewClient(c.RiotAPIKey)

	datasource := internal.NewDatasource(client, store)

	server := frontend.New(&frontend.Handler{Datasource: datasource})
	server.Logger = logger

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		err := http.ListenAndServe(address, server)
		logger.Error("error starting server", "err", err)
	}()

	logger.Info("start server", "addr", address)

	<-ctx.Done()

	logger.Info("stop server")

	return nil
}

func connectPostgres(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
