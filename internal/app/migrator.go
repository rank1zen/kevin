package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/rank1zen/kevin/internal/config"
	"github.com/rank1zen/kevin/internal/log"
	"github.com/rank1zen/kevin/migrations"
)

type Migrator struct {
	config       *config.Config
	logger       *slog.Logger
	postgresConn *pgx.Conn
	migrator     *migrate.Migrator

	errors []error
}

func NewMigrator(ctx context.Context) *Migrator {
	app := &Migrator{}
	app.logger = log.New()

	cfg, err := config.NewConfig()
	if err != nil {
		app.errors = append(app.errors, fmt.Errorf("failed to load configuration: %w", err))
		return app
	}
	app.config = cfg

	conn, err := pgx.Connect(ctx, cfg.GetDatabaseURL())
	if err != nil {
		app.errors = append(app.errors, fmt.Errorf("failed to connect to postgres: %w", err))
		return app
	}
	app.postgresConn = conn

	m, err := migrate.NewMigrator(ctx, app.postgresConn, "public.schema_version")
	if err != nil {
		app.errors = append(app.errors, fmt.Errorf("failed to create migrator: %w", err))
		return app
	}

	if err := m.LoadMigrations(migrations.Migrations); err != nil {
		app.errors = append(app.errors, fmt.Errorf("failed to load migrations: %w", err))
		return app
	}
	app.migrator = m

	return app
}

func (a *Migrator) Run(ctx context.Context) int {
	if len(a.errors) > 0 {
		a.logger.Error("app startup failed", "err", errors.Join(a.errors...))
		return 1
	}

	a.logger.Info("running migrations")

	if err := a.migrator.Migrate(ctx); err != nil {
		a.logger.Error("migration failed", "err", err)
		return 1
	}

	if a.postgresConn != nil {
		_ = a.postgresConn.Close(ctx)
	}

	return 0
}

func (a *Migrator) Errors() error {
	return errors.Join(a.errors...)
}
