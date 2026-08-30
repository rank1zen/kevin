package pgtestcontainer

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/tern/v2/migrate"
	"github.com/rank1zen/kevin/migrations"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PGInstance struct {
	container *pg.PostgresContainer
	pool      *pgxpool.Pool
	pgURL     string
}

// NewPGInstance sets up a Postgres server in a docker container. It will use
// the current schema version.
func NewPGInstance(ctx context.Context) *PGInstance {
	const (
		pgDBName   = "postgres_test"
		pgUser     = "kevin"
		pgPassword = "secret"
		pgImage    = "docker.io/postgres:18-alpine"
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
	if err != nil {
		log.Fatal(err)
	}

	pgInstance := &PGInstance{
		container: container,
		pgURL:     pgURL,
	}

	pgInstance.migrateSchema(ctx)

	pool, err := pgxpool.New(ctx, pgInstance.pgURL)
	if err != nil {
		log.Fatal(err)
	}
	pgInstance.pool = pool

	return pgInstance
}

// SetupTx starts a new database transaction for a test. It is safe to be used in parallel.
func (p *PGInstance) SetupTx(t testing.TB) pgx.Tx {
	tx, err := p.pool.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := tx.Rollback(context.Background())
		if err != nil {
			t.Fatal(err)
		}
	})

	return tx
}

func (p *PGInstance) migrateSchema(ctx context.Context) {
	conn, err := pgx.Connect(ctx, p.pgURL)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = conn.Close(ctx)
	}()

	m, err := migrate.NewMigrator(ctx, conn, "public.schema_version")
	if err != nil {
		log.Fatal(err)
	}

	if err := m.LoadMigrations(migrations.Migrations); err != nil {
		log.Fatal(err)
	}

	if err = m.Migrate(ctx); err != nil {
		log.Fatal(err)
	}
}

func (p *PGInstance) Terminate(ctx context.Context) error {
	return p.container.Terminate(ctx)
}

func (p *PGInstance) GetConnectionString() string {
	return p.pgURL
}
