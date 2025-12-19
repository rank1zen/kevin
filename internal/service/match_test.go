package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchService_GetMatchDetail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	ds := SetupDatasource(ctx, t)

	req := service.GetMatchDetailRequest{
		MatchID: "NA1_5346312088",
	}

	match, err := (*service.MatchService)(ds).GetMatchDetail(ctx, req)
	require.NoError(t, err)

	T1 := findT1(t, *match)

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct version",
			Expected: "15.16.702.7993",
			Actual:   match.Version,
		},
		{
			Name:     "expects correct summoner spells",
			Expected: [2]int{4, 12},
			Actual:   T1.SummonerIDs,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}

	assert.Equal(t, riot.RegionNA1, *req.Region)
	assert.EqualValues(t, T1.PUUID, "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")
}

func TestMatchService_GetMatchlist(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	ds := SetupDatasource(ctx, t)

	req := service.GetMatchlistRequest{
		PUUID: "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg",
	}
	req.StartTS = new(time.Time)
	*req.StartTS = time.Date(2025, 6, 10, 10, 59, 0, 0, time.UTC)

	req.EndTS = new(time.Time)
	*req.EndTS = time.Date(2025, 6, 10, 11, 0, 0, 0, time.UTC)

	storeMatches, err := (*service.MatchService)(ds).GetMatchlist(ctx, req)
	require.NoError(t, err)

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct match count",
			Expected: 1,
			Actual:   len(storeMatches),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}
