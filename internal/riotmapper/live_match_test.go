package riotmapper_test

import (
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal/riotmapper"
	"github.com/rank1zen/kevin/internal/sample"
	"github.com/stretchr/testify/assert"
)

func TestMapLiveMatch(t *testing.T) {
	riotLiveMatch := sample.WithSampleLiveMatch()

	actual := riotmapper.MapLiveMatch(&riotLiveMatch)

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

	actualParticipant := actual.Participants[0]

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct participant champion",
			Expected: 517,
			Actual:   actualParticipant.ChampionID,
		},
		{
			Name:     "expects correct participant team id",
			Expected: 100,
			Actual:   actualParticipant.TeamID,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}
