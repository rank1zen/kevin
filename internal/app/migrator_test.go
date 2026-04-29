package app_test

import (
	"context"
	"testing"

	"github.com/rank1zen/kevin/internal/app"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestMigrator_Run(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:18",
		postgres.WithDatabase("kevin_test"),
		postgres.WithUsername("kevin"),
		postgres.WithPassword("kevin"),
		postgres.BasicWaitStrategies(),
	)

	require.NoError(t, err)

	t.Cleanup(func() {
		_ = pgContainer.Terminate(ctx)
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	t.Setenv("KEVIN_DATABASE_URL", connStr)
	t.Setenv("KEVIN_RIOT_API_KEY", "test-key")
	t.Setenv("KEVIN_ENV", "production")
	t.Setenv("PORT", "4099")

	a := app.NewMigrator(ctx)
	require.Empty(t, a.Errors(), "app.New() should not have errors")

	code := a.Run(ctx)
	require.Equal(t, 0, code)
}
