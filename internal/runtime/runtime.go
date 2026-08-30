package runtime

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal/riot"
)

type Runtime struct {
	config *Config

	logger *slog.Logger

	pgxPool    *pgxpool.Pool
	riotClient *riot.Client

	httpServer *http.Server
}

func NewRuntime(ctx context.Context, config *Config) (*Runtime, error) {
	return startUpRuntime(ctx, config)
}

func (r *Runtime) Run(ctx context.Context) int {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	serverErrCh := make(chan error, 1)
	go func() {
		err := r.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
	}()

	r.logger.Info("server started", "port", r.config.Port, "environment", r.config.Environment)

	select {
	case err := <-serverErrCh:
		r.logger.Error("listen and serve", "err", err)
		r.tearDown()
		return 1
	case <-ctx.Done():
	}

	r.logger.Info("shutting down server")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := r.httpServer.Shutdown(shutdownCtx); err != nil {
		r.logger.Error("server shutdown failed", "err", err)
	}

	r.tearDown()
	return 0
}

// startUpRuntime initializes all runtime dependencies from the provided config.
func startUpRuntime(ctx context.Context, config *Config) (*Runtime, error) {
	logger := slog.Default()

	err := validateConfigForRuntime(config)
	if err != nil {
		return nil, err
	}

	riotClient := riot.NewClient(config.RiotAPIKey)

	pgxPool, err := initializePostgres(ctx, *config)
	if err != nil {
		return nil, err
	}

	httpServer := initializeServer(*config)

	return &Runtime{
		config:     config,
		logger:     logger,
		pgxPool:    pgxPool,
		riotClient: riotClient,
		httpServer: httpServer,
	}, nil
}

func (r *Runtime) tearDown() {
	if r.pgxPool != nil {
		r.pgxPool.Close()
	}
}

func validateConfigForRuntime(cfg *Config) error {
	var errs []error

	if cfg.RiotAPIKey == "" {
		errs = append(errs, errors.New("KEVIN_RIOT_API_KEY is not set"))
	}

	if cfg.DatabaseURL == "" {
		errs = append(errs, errors.New("database url is required"))
	}

	_, err := pgx.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
