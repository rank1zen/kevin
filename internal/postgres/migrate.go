package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/rank1zen/kevin/migrations"
)

func Migrate(ctx context.Context, connString string) error {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return err
	}

	defer func() {
		_ = conn.Close(ctx)
	}()

	m, err := migrate.NewMigrator(ctx, conn, "public.schema_version")
	if err != nil {
		return err
	}

	if err := m.LoadMigrations(migrations.Migrations); err != nil {
		return err
	}

	if err = m.Migrate(ctx); err != nil {
		return err
	}

	return err
}
