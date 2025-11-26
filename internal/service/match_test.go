package service_test

import (
	"context"
	"testing"

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
