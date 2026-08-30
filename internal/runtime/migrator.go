package runtime

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/rank1zen/kevin/migrations"
)

func initializeMigrator(ctx context.Context, cfg Config, pgxConn *pgx.Conn) (*migrate.Migrator, error) {
	m, err := migrate.NewMigrator(ctx, pgxConn, "public.schema_version")
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	err = m.LoadMigrations(migrations.Migrations)
	if err != nil {
		return nil, fmt.Errorf("failed to load migrations: %w", err)
	}

	return m, nil
}
