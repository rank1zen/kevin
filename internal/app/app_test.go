package app_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal/app"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestAppStartsAndResponds(t *testing.T) {
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

	a := app.New(ctx)
	require.Empty(t, a.Errors(), "app.New() should not have errors")

	go a.Run(ctx)

	require.Eventually(t, func() bool {
		resp, err := http.Get("http://localhost:4099/healthz")
		if err != nil {
			return false
		}
		_ = resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	},
		10*time.Second,
		100*time.Millisecond,
		"server never became ready",
	)

	resp, err := http.Get("http://localhost:4099/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
