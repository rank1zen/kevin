package internal_test

import (
	"context"
	"testing"

	"github.com/rank1zen/kevin/internal"
	"github.com/stretchr/testify/assert"
)

func TestMatchService_GetMatchDetail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	ds := SetupDatasource(ctx, t)
	service := (*internal.MatchService)(ds)

	req := internal.GetMatchDetailRequest{
		MatchID: "NA1_5346312088",
	}

	match, err := service.GetMatchDetail(ctx, req)
	if assert.NoError(t, err) {
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
	}
}
