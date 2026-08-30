package runtime_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuntime_Run(t *testing.T) {
	t.Setenv("KEVIN_DATABASE_URL", DefaultPGInstance.GetConnectionString())

	t.Run("should fail if riot api key is not found", func(t *testing.T) {
		cfg, err := runtime.NewConfig()
		require.NoError(t, err)

		_, err = runtime.NewRuntime(t.Context(), cfg)
		assert.Error(t, err)
	})

	t.Run("should be healthy if started", func(t *testing.T) {
		t.Setenv("KEVIN_RIOT_API_KEY", "test-key")

		config, err := runtime.NewConfig()
		require.NoError(t, err)

		rt, err := runtime.NewRuntime(t.Context(), config)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()

		codeCh := make(chan int, 1)
		go func() {
			codeCh <- rt.Run(ctx)
		}()

		require.Eventually(
			t,
			func() bool {
				res, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz", config.Port))
				if err != nil {
					return false
				}
				defer func() {
					_ = res.Body.Close()
				}()

				return res.StatusCode == http.StatusOK
			},
			3*time.Second,
			50*time.Millisecond,
			"server never became healthy",
		)

		cancel()
		require.EqualValues(t, 0, <-codeCh)
	})
}
