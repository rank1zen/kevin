package runtime_test

import (
	"testing"

	"github.com/rank1zen/kevin/internal/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Setenv("KEVIN_RIOT_API_KEY", "test-key")
	t.Setenv("KEVIN_ENV", "dev")
	t.Setenv("PORT", "8080")

	t.Run("should bind envs to struct", func(t *testing.T) {
		cfg, err := runtime.NewConfig()
		require.NoError(t, err)

		assert.Equal(t, "test-key", cfg.RiotAPIKey)
		assert.Equal(t, "dev", cfg.Environment)
		assert.Equal(t, 8080, cfg.Port)
	})
}
