package riot_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLeagueGetLeagueEntriesByPUUID(t *testing.T) {
	ctx := context.Background()

	client, server := MakeTestClient(t, http.StatusOK, "sample/lol/league/v4/entries/by-puuid/0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q.json",
		func(r *http.Request) {
			for _, tc := range []struct {
				Name             string
				Expected, Actual any
			}{
				{
					Name:     "expects correct endpoint",
					Expected: "/lol/league/v4/entries/by-puuid/0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
					Actual:   r.URL.Path,
				},
			} {
				t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
			}
		},
	)

	defer server.Close()

	entries, err := client.League.GetLeagueEntriesByPUUID(ctx, riot.RegionNA1, OrrangePUUID.String())

	require.NoError(t, err)

	require.Len(t, entries, 1)
	soloq := entries[0]

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct puuid for orrange",
			Expected: riot.QueueTypeRankedSolo5x5,
			Actual:   soloq.QueueType,
		},
		{
			Name:     "expects correct name for orrange",
			Expected: riot.TierEmerald,
			Actual:   soloq.Tier,
		},
		{
			Name:     "expects correct name for orrange",
			Expected: riot.Division4,
			Actual:   soloq.Division,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}
