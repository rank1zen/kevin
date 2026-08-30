package runtime_test

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/kevin/internal/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrationRuntime_Run(t *testing.T) {
	t.Setenv("KEVIN_DATABASE_URL", DefaultPGInstance.GetConnectionString())

	cfg, err := runtime.NewConfig()
	require.NoError(t, err)

	t.Run("should apply migrations to clean db", func(t *testing.T) {

		rt, err := runtime.NewMigrationRuntime(t.Context(), cfg)
		require.NoError(t, err)

		code := rt.Run(t.Context())
		require.Equal(t, 0, code)

		pgxConn, err := pgx.Connect(t.Context(), DefaultPGInstance.GetConnectionString())
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, pgxConn.Close(t.Context()))
		})

		var count int
		err = pgxConn.QueryRow(t.Context(), "select count(*) from public.schema_version").
			Scan(&count)
		require.NoError(t, err)
		assert.Greater(t, count, 0)
	})
}
