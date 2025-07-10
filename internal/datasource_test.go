package internal_test

import (
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStartAndEndUnix(t *testing.T) {
	t.Run(
		"timestamps in Toronto",
		func(t *testing.T) {
			location, err := time.LoadLocation("America/Toronto")
			require.NoError(t ,err)

			local := time.Date(2025, 7, 9, 19, 11, 40, 0, location)
			start, end := internal.GetStartAndEndUnix(local)

			expectedStart := time.Date(2025, 7, 9, 0, 0, 0, 0, location)

			assert.Equal(t, expectedStart.Unix(), start)
			print(start)
			assert.Equal(t, expectedStart.Add(24*time.Hour - 1*time.Second).Unix(), end)
		},
	)
}
