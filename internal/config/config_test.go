package config_test

import (
	"testing"

	"github.com/rank1zen/kevin/internal/config"
	"github.com/stretchr/testify/require"
)

func setBaseEnv(t *testing.T) {
	t.Helper()

	t.Setenv("KEVIN_RIOT_API_KEY", "test-key")
	t.Setenv("KEVIN_ENV", "dev")
	t.Setenv("PORT", "8080")
}

func clearEnv(t *testing.T) {
	t.Helper()

	t.Setenv("KEVIN_DATABASE_URL", "")
	t.Setenv("KEVIN_RIOT_API_KEY", "")
	t.Setenv("KEVIN_ENV", "")
	t.Setenv("PORT", "")
}

func TestNewConfig_ValuesAreSet(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("KEVIN_DATABASE_URL", "postgres://user:pass@localhost:5432/kevin")

	cfg, err := config.NewConfig()
	require.NoError(t, err)

	require.Equal(t, "test-key", cfg.RiotAPIKey)
	require.Equal(t, config.Development, cfg.Environment)
	require.Equal(t, 8080, cfg.Port)
	require.Equal(t, "postgres://user:pass@localhost:5432/kevin", cfg.DatabaseURL)
}
func TestNewConfig_DefaultValuesAreSet(t *testing.T) {
	clearEnv(t)
	t.Setenv("KEVIN_DATABASE_URL", "postgres://user:pass@localhost:5432/kevin")
	t.Setenv("KEVIN_RIOT_API_KEY", "test-key")

	cfg, err := config.NewConfig()
	require.NoError(t, err)

	require.Equal(t, config.Development, cfg.Environment)
	require.Equal(t, 7331, cfg.Port)
}
func TestReadsTwoPortVariables(t *testing.T) {
	t.Setenv("KEVIN_DATABASE_URL", "postgres://user:pass@localhost:5432/kevin")
	t.Setenv("KEVIN_RIOT_API_KEY", "test-key")

	t.Setenv("KEVIN_PORT", "8080")
	t.Setenv("PORT", "8090")

	cfg, err := config.NewConfig()
	require.NoError(t, err)
	require.Equal(t, 8090, cfg.Port)

	t.Setenv("PORT", "")
	cfg, err = config.NewConfig()
	require.NoError(t, err)
	require.Equal(t, 8080, cfg.Port)
}
