package migrate_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal/migrate"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	dbname = "internal-migrate"
	dbuser = "yujin"
	dbpass = "secret"
)

var (
	container *postgres.PostgresContainer
	conn      *pgx.Conn
	dburl     string
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

	conn, err = pgx.Connect(ctx, dburl)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	err = container.Terminate(ctx)
	if err != nil {
		log.Fatalf("terminating: %s", err)
	}

	os.Exit(code)
}

func TestMigrate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	err := migrate.Migrate(ctx, conn)
	assert.NoError(t, err)
}
