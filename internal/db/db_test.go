package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rank1zen/yujin/internal/db/migrate"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	container *postgres.PostgresContainer
	dburl     string
)

const (
	dbname = "testing"
	dbuser = "yujin"
	dbpass = "secret"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	container, err = postgres.Run(ctx, "docker.io/postgres:16-alpine",
		postgres.WithDatabase(dbname),
		postgres.WithUsername(dbuser),
		postgres.WithPassword(dbpass),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"))
	if err != nil {
		log.Fatalf("running postgres container: %s", err)
	}

	dburl, err = container.ConnectionString(ctx)
	if err != nil {
		log.Fatal(err)
	}

	pgxconn, err := pgx.Connect(ctx, dburl)
	if err != nil {
		log.Fatalf("migration conn: %s", err)
	}

	err = migrate.Migrate(ctx, pgxconn)
	if err != nil {
		log.Fatalf("migrating: %s", err)
	}

	pgxconn.Close(ctx)

	err = container.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		log.Fatalf("creating snapshot: %s", err)
	}

	code := m.Run()

	err = container.Terminate(ctx)
	if err != nil {
		log.Fatalf("terminating: %s", err)
	}

	os.Exit(code)
}

func setupDB(t testing.TB) *DB {
	ctx := context.Background()

	db, err := NewDB(ctx, dburl)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()

		err := container.Restore(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	return db
}
