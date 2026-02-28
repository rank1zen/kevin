package config_test

import (
	"os"
	"testing"

	"github.com/rank1zen/kevin/internal/config"
	"github.com/stretchr/testify/require"
)

func setBaseEnv(t *testing.T) {
	t.Helper()

	t.Setenv("KEVIN_RIOT_API_KEY", "test-key")
	t.Setenv("KEVIN_ENV", "development")
	t.Setenv("PORT", "8080")
}

func TestNewConfig_ValuesAreSet(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("KEVIN_DATABASE_URL", "postgres://user:pass@localhost:5432/kevin")

	cfg, err := config.NewConfig()
	require.NoError(t, err)

	require.Equal(t, "test-key", cfg.GetRiotAPIKey())
	require.True(t, cfg.IsDevelopment())
	require.Equal(t, 8080, cfg.GetPort())
	require.Equal(t, "postgres://user:pass@localhost:5432/kevin", cfg.GetDatabaseURL())
}

func TestNewConfig_RejectsWhenMissingAny(t *testing.T) {
	os.Clearenv()
	_, err := config.NewConfig()
	require.Error(t, err)
}

func TestNewConfig_RejectsDatabaseURLWithInvalidScheme(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("KEVIN_DATABASE_URL", "mysql://user:pass@localhost:3306/kevin")

	_, err := config.NewConfig()
	require.Error(t, err)
	require.ErrorContains(t, err, `KEVIN_DATABASE_URL: url must use postgres/postgresql scheme`)
}

func TestNewConfig_RejectsDatabaseURLWithoutHost(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("KEVIN_DATABASE_URL", "postgres:///kevin")

	_, err := config.NewConfig()
	require.Error(t, err)
	require.ErrorContains(t, err, "KEVIN_DATABASE_URL: url must include a host")
}
