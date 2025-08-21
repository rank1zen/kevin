package riot_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchGetMatchList(t *testing.T) {
	ctx := context.Background()

	client, server := MakeTestClient(t, http.StatusOK, "sample/lol/match/v5/matches/by-puuid/0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q/ids.json",
		func(r *http.Request) {
			for _, tc := range []struct {
				Name             string
				Expected, Actual any
			}{
				{
					Name:     "expects correct endpoint",
					Expected: "/lol/match/v5/matches/by-puuid/0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q/ids",
					Actual:   r.URL.Path,
				},
				{
					Name:     "expects correct start time query",
					Expected: "1751601600",
					Actual:   r.URL.Query().Get("startTime"),
				},
				{
					Name:     "expects correct end time query",
					Expected: "1751687999",
					Actual:   r.URL.Query().Get("endTime"),
				},
				{
					Name:     "expects correct queue query",
					Expected: "420",
					Actual:   r.URL.Query().Get("queue"),
				},
			} {
				t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
			}
		},
	)

	defer server.Close()

	options := riot.MatchListOptions{
		StartTime: new(int64),
		EndTime:   new(int64),
		Queue:     new(int),
		Start:     0,
		Count:     20,
	}

	*options.StartTime = 1751601600
	*options.EndTime = 1751687999
	*options.Queue = 420

	matches, err := client.Match.GetMatchList(ctx, riot.RegionNA1, OrrangePUUID.String(), options)
	require.NoError(t, err)

	assert.Len(t, matches, 20)
}

func TestMatchGetMatch(t *testing.T) {
	ctx := context.Background()

	client, server := MakeTestClient(t, http.StatusOK, "sample/lol/match/v5/matches/NA1_5304757838.json",
		func(r *http.Request) {
			for _, tc := range []struct {
				Name             string
				Expected, Actual any
			}{
				{
					Name:     "expects correct endpoint",
					Expected: "/lol/match/v5/matches/NA1_5304757838",
					Actual:   r.URL.Path,
				},
			} {
				t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
			}
		},
	)

	defer server.Close()

	match, err := client.Match.GetMatch(ctx, riot.RegionNA1, "NA1_5304757838")
	require.NoError(t, err)

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects match id",
			Expected: "NA1_5304757838",
			Actual:   match.Metadata.MatchID,
		},
		{
			Name:     "expects correct kills for participant 0",
			Expected: 8,
			Actual:   match.Info.Participants[0].Kills,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}
