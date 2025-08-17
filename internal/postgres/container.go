package postgres

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/tern/v2/migrate"
	"github.com/rank1zen/kevin/internal"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PGInstance struct {
	container *pg.PostgresContainer

	pgURL string
}

// NewPGInstance sets up a postgres server in a docker container. It will use
// the current schema version.
func NewPGInstance(ctx context.Context) *PGInstance {
	const (
		pgDBName   = "postgres_test"
		pgUser     = "kevin"
		pgPassword = "secret"
		pgImage    = "docker.io/postgres:16-alpine"
	)

	container, err := pg.Run(ctx, pgImage,
		pg.WithDatabase(pgDBName),
		pg.WithUsername(pgUser),
		pg.WithPassword(pgPassword),
		pg.BasicWaitStrategies(),
		pg.WithSQLDriver("pgx"),
	)

	if err != nil {
		log.Fatalf("running postgres container: %s", err)
	}

	pgURL, err := container.ConnectionString(ctx)

	pgInstance := &PGInstance{
		container,
		pgURL,
	}

	pgInstance.migrateSchema(ctx)

	if err := pgInstance.container.Snapshot(ctx, pg.WithSnapshotName("test-snapshot")); err != nil {
		log.Fatalf("creating snapshot: %s", err)
	}

	return pgInstance
}

// SetupStore creates an empty store, and will clean up after test t finishes.
func (p *PGInstance) SetupStore(ctx context.Context, t testing.TB) internal.Store {
	conn, err := pgxpool.New(ctx, p.pgURL)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		conn.Close()

		if err := p.container.Restore(ctx); err != nil {
			t.Fatal(err)
		}
	})

	store := NewStore(conn)

	return store
}

func (p *PGInstance) SetupConn(ctx context.Context, t testing.TB) *pgxpool.Pool {
	conn, err := pgxpool.New(ctx, p.pgURL)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		conn.Close()

		if err := p.container.Restore(ctx); err != nil {
			t.Fatal(err)
		}
	})

	return conn
}

func (p *PGInstance) migrateSchema(ctx context.Context) {
	conn, err := pgx.Connect(ctx, p.pgURL)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close(ctx)

	m, err := migrate.NewMigrator(ctx, conn, "public.schema_version")
	if err != nil {
		log.Fatal(err)
	}

	m.LoadMigrations(os.DirFS("../../migrations"))

	if err = m.Migrate(ctx); err != nil {
		log.Fatal(err)
	}
}

func (p *PGInstance) Terminate(ctx context.Context) error {
	return p.container.Terminate(ctx)
}
