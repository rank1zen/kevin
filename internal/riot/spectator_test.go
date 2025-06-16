package riot_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpectatorGetLiveMatch(t *testing.T) {
	ctx := context.Background()

	client, server := MakeTestClient(t, http.StatusOK, "sample/lol/spectator/v5/active-games/by-summoner/_goIttvCedN8i6rEm3no4S6ZgSRfn4YM-VMl7IQtVdulyr2je8OMYB9_15o067LHjsLZcOEwixkfTA.json",
		func(r *http.Request) {
			for _, tc := range []struct {
				Name             string
				Expected, Actual any
			}{
				{
					Name:     "expects correct endpoint",
					Expected: "/lol/spectator/v5/active-games/by-summoner/_goIttvCedN8i6rEm3no4S6ZgSRfn4YM-VMl7IQtVdulyr2je8OMYB9_15o067LHjsLZcOEwixkfTA",
					Actual:   r.URL.Path,
				},
			} {
				t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
			}
		},
	)

	defer server.Close()

	match, err := client.Spectator.GetLiveMatch(ctx, riot.RegionNA1, "_goIttvCedN8i6rEm3no4S6ZgSRfn4YM-VMl7IQtVdulyr2je8OMYB9_15o067LHjsLZcOEwixkfTA")
	require.NoError(t, err)

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects match id",
			Expected: 5330985291,
			Actual:   match.GameID,
		},
		{
			Name:     "expects correct puuid for participant 0",
			Expected: "DBorWJClDRI17ivbgFYmR6QEdH7J6GTeGGbw4YWc6JZtytIX10kUwD6_ZmJkkQ2gchUkWlxOAUMkIg",
			Actual:   match.Participants[0].PUUID,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}

func TestSpectatorGetLiveMatch_NotFound(t *testing.T) {
	ctx := context.Background()

	client, server := MakeTestClient(t, http.StatusNotFound, "sample/lol/spectator/v5/active-games/by-summoner/0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q.json")
	defer server.Close()

	_, err := client.Spectator.GetLiveMatch(ctx, riot.RegionNA1, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q")
	require.ErrorIs(t, err, riot.ErrNotFound)
}
