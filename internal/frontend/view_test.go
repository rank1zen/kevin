package frontend_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/component"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/sample"
	"github.com/stretchr/testify/assert"
)

func TestLiveMatch(t *testing.T) {
	ctx := context.Background()
	v := frontend.LiveMatch{ }

	c := v.ToTempl(ctx)

	// assert expectation
}

func TestNewMatchDetail(t *testing.T) {
	actual := frontend.NewMatchDetail(internal.NewMatch(sample.WithSampleMatch()))

	actualPlayer := actual.BlueSide[0]

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct date",
			Expected: time.UnixMilli(1749596377340),
			Actual:   actual.Date,
		},
		{
			Name: "expects correct kda widget",
			Expected: component.KDAWidget{
				Kills:             8,
				Deaths:            0,
				Assists:           1,
				KillParticipation: 0,
				KilLDeathRatio:    9,
			},
			Actual: actualPlayer.KDAWidget,
		},
		{
			Name:     "expects correct rank delta widget",
			Expected: component.RankDeltaWidget{RankWidget: nil, LPChange: nil, Win: true},
			Actual:   actualPlayer.RankDeltaWidget,
		},
		{
			Name:     "expects bot",
			Expected: "ADADADA",
			Actual:   actual.BlueSide[3].PUUID,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}

func TestNewMatchHistory(t *testing.T) {}

func TestNewSearchResult(t *testing.T) {
}
