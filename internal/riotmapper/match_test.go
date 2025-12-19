package riotmapper_test

import (
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riotmapper"
	"github.com/rank1zen/kevin/internal/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapMatch(t *testing.T) {
	riotMatch := sample.WithSampleMatch()

	actualMatch := riotmapper.MapMatch(riotMatch)

	for _, tc := range []struct {
		Name     string
		Expected any
		Actual   any
	}{
		{
			Name:     "expects correct ID",
			Expected: "NA1_5304757838",
			Actual:   actualMatch.ID,
		},
		{
			Name:     "expects correct date",
			Expected: time.UnixMilli(1749596377340),
			Actual:   actualMatch.Date,
		},
		{
			Name:     "expects correct duration",
			Expected: 1131 * time.Second,
			Actual:   actualMatch.Duration,
		},
		{
			Name:     "expects correct version",
			Expected: "15.11.685.5259",
			Actual:   actualMatch.Version,
		},
		{
			Name:     "expects correct winner",
			Expected: 100,
			Actual:   actualMatch.WinnerID,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}

	var actualParticipant *internal.Participant
	for _, p := range actualMatch.Participants {
		if p.PUUID == "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg" {
			actualParticipant = &p
		}
	}

	require.NotNil(t, actualParticipant)

	for _, tc := range []struct {
		Name     string
		Expected any
		Actual   any
	}{
		{
			Name:     "expects correct participant champion",
			Expected: 63,
			Actual:   actualParticipant.ChampionID,
		},
		{
			Name:     "expects correct participant kills",
			Expected: 2,
			Actual:   actualParticipant.Kills,
		},
		{
			Name:     "expects correct participant kill participation",
			Expected: float32(10.0 / 27.0),
			Actual:   actualParticipant.KillParticipation,
		},
		{
			Name:     "expects correct participant cs per minute",
			Expected: float32(131.0 * 60 / 1131),
			Actual:   actualParticipant.CreepScorePerMinute,
		},
		{
			Name:     "expects correct participant damage percentage",
			Expected: float32(12629.0 / 56169),
			Actual:   actualParticipant.DamagePercentageTeam,
		},
		{
			Name:     "expects correct participant gold percentage",
			Expected: float32(6856.0 / 41017),
			Actual:   actualParticipant.GoldPercentageTeam,
		},
		{
			Name:     "expects correct participant position",
			Expected: internal.TeamPositionMiddle,
			Actual:   actualParticipant.TeamPosition,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}
