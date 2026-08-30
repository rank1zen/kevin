package runtime

import (
	"context"
	"errors"
	"fmt"
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
func startUpMigrationRuntime(ctx context.Context, cfg *Config) (*MigrationRuntime, error) {
	logger := slog.Default()
	logger.Info("starting migration runtime")

	err := validateConfigForMigrationRuntime(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	pgxConn, err := initializePgxConn(ctx, *cfg)
	if err != nil {
		return nil, err
	}

	migrator, err := initializeMigrator(ctx, *cfg, pgxConn)
	if err != nil {
		return nil, err
	}

	logger.Info("migration runtime initialized successfully")
	return &MigrationRuntime{
		config:   cfg,
		logger:   logger,
		pgxConn:  pgxConn,
		migrator: migrator,
	}, nil
}

func validateConfigForMigrationRuntime(cfg *Config) error {
	var errs []error

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
