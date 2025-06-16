package internal_test

import (
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMatch(t *testing.T) {
	actual := internal.NewMatch(sample.WithSampleMatch())

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct ID",
			Expected: "NA1_5304757838",
			Actual:   actual.ID,
		},
		{
			Name:     "expects correct date",
			Expected: time.UnixMilli(1749596377340),
			Actual:   actual.Date,
		},
		{
			Name:     "expects correct duration",
			Expected: 1131 * time.Second,
			Actual:   actual.Duration,
		},
		{
			Name:     "expects correct version",
			Expected: "15.11.685.5259",
			Actual:   actual.Version,
		},
		{
			Name:     "expects correct winner",
			Expected: 100,
			Actual:   actual.WinnerID,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}

	var t1 *internal.Participant

	for _, p := range actual.Participants {
		if p.PUUID == sample.SummonerT1OKGOODYESNA1.PUUID {
			t1 = &p
		}
	}

	require.NotNil(t, t1)
	require.Equal(t, actual.ID, t1.MatchID)

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct participant champion",
			Expected: 63,
			Actual:   t1.ChampionID,
		},
		{
			Name:     "expects correct participant kills",
			Expected: 2,
			Actual:   t1.Kills,
		},
		{
			Name:     "expects correct participant kill participation",
			Expected: float32(10.0 / 27.0),
			Actual:   t1.KillParticipation,
		},
		{
			Name:     "expects correct participant cs per minute",
			Expected: float32(131.0 * 60 / 1131),
			Actual:   t1.CreepScorePerMinute,
		},
		{
			Name:     "expects correct participant damage percentage",
			Expected: float32(12629.0 / 56169),
			Actual:   t1.DamagePercentageTeam,
		},
		{
			Name:     "expects correct participant gold percentage",
			Expected: float32(6856.0 / 41017),
			Actual:   t1.GoldPercentageTeam,
		},
		{
			Name:     "expects correct participant position",
			Expected: internal.TeamPositionMiddle,
			Actual:   t1.TeamPosition,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}

func TestNewLiveMatch(t *testing.T) {
	actual := internal.NewLiveMatch(sample.WithSampleLiveMatch())

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct ID",
			Expected: "NA1_5330985291",
			Actual:   actual.ID,
		},
		{
			Name:     "expects correct date",
			Expected: time.UnixMilli(1753291144171),
			Actual:   actual.Date,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}

	t1 := actual.Participants[0]
	require.Equal(t, actual.ID, t1.MatchID)

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct participant champion",
			Expected: 517,
			Actual:   t1.ChampionID,
		},
		{
			Name:     "expects correct participant team id",
			Expected: 100,
			Actual:   t1.TeamID,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}
