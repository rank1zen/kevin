// app is responsible for the runtime.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal/config"
	"github.com/rank1zen/kevin/internal/log"
	"github.com/rank1zen/kevin/internal/profile"
	dbProfile "github.com/rank1zen/kevin/internal/profile/db"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/route"
)

type App struct {
	config       *config.Config
	logger       *slog.Logger
	postgresConn *pgxpool.Pool
	riotClient   *riot.Client
	server       http.Handler

	errors []error
}

func New(ctx context.Context) *App {
	app := &App{}
	app.logger = log.New()

	cfg, err := config.NewConfig()
	if err != nil {
		app.errors = append(app.errors, fmt.Errorf("failed to load configuration: %w", err))
		return app
	}
	app.config = cfg

	pool, err := connectPostgres(ctx, cfg.GetDatabaseURL())
	if err != nil {
		app.errors = append(app.errors, fmt.Errorf("failed to connect to postgres: %w", err))
		return app
	}
	app.postgresConn = pool

	riotClient := riot.NewClient(cfg.GetRiotAPIKey())
	app.riotClient = riotClient

	app.server = route.Router(
		riotClient,
		profile.NewProfileService(riotClient, dbProfile.NewStore(pool)),
	)

	return app
}

func (a *App) Run(ctx context.Context) int {
	if len(a.errors) > 0 {
		a.logger.Error("app startup failed", "err", errors.Join(a.errors...))
		return 1
	}

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.GetPort()),
		Handler: a.server,
	}

	serverErrCh := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
	}()

	a.logger.Info("server started", "address", a.config.GetPort(), "environment", a.config.IsDevelopment())

	select {
	case err := <-serverErrCh:
		a.logger.Error("listen and serve", "err", err)
		if a.postgresConn != nil {
			a.postgresConn.Close()
		}
		return 1
	case <-ctx.Done():
	}

	a.logger.Info("shutting down server")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("server shutdown failed", "err", err)
	}

	if a.postgresConn != nil {
		a.postgresConn.Close()
	}

	return 0
}

func (a *App) Errors() error {
	return errors.Join(a.errors...)
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

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}
