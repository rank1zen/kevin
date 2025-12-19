package internal_test

import (
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/sample"
	"github.com/stretchr/testify/assert"
)

func TestRiotToProfileMapper(t *testing.T) {
	riotAccount := sample.Account(t)
	riotLeagueList := sample.LeagueList(t)

	mapper := internal.RiotToProfileMapper{
		Account:       riotAccount,
		Rank:          &riotLeagueList[0],
		EffectiveDate: time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC),
	}

	actual := mapper.Convert()

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct puuid",
			Expected: riot.PUUID("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg"),
			Actual:   actual.PUUID,
		},
		{
			Name:     "expects correct effective date",
			Expected: time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC),
			Actual:   actual.Rank.EffectiveDate,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}
