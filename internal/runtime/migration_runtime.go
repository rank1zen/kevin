package runtime

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
)

type MigrationRuntime struct {
	config *Config

	logger *slog.Logger

	pgxConn *pgx.Conn

	migrator *migrate.Migrator
}

func NewMigrationRuntime(ctx context.Context, config *Config) (*MigrationRuntime, error) {
	return startUpMigrationRuntime(ctx, config)
}

func (r *MigrationRuntime) Run(ctx context.Context) int {
	r.logger.Info("starting migration")

	err := r.migrator.Migrate(ctx)
	if err != nil {
		r.logger.Error("migration failed", "err", err)
		return 1
	}

	r.tearDown(ctx)
	r.logger.Info("migration complete")
	return 0
}

func (r *MigrationRuntime) tearDown(ctx context.Context) {
	if r.pgxConn != nil {
		_ = r.pgxConn.Close(ctx)
	}
}

// startU initializes all runtime dependencies from the provided config.
func startUpMigrationRuntime(ctx context.Context, config *Config) (*MigrationRuntime, error) {
	logger := slog.Default()

	pgxConn, err := initializePgxConn(ctx, *config)
	if err != nil {
		return nil, err
	}

	migrator, err := initializeMigrator(ctx, *config, pgxConn)
	if err != nil {
		return nil, err
	}

	return &MigrationRuntime{
		config:   config,
		logger:   logger,
		pgxConn:  pgxConn,
		migrator: migrator,
	}, nil
}
