package internal_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRiotToMatchMapper(t *testing.T) {
	riotMatch := sample.WithSampleMatch()

	mapper := internal.RiotToMatchMapper{
		Match: riotMatch,
	}

	actual := mapper.Map()

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
}

func TestRiotToParticipantMapper(t *testing.T) {
	riotMatch := sample.WithSampleMatch()

	var riotParticipant *riot.MatchParticipant
	for _, p := range riotMatch.Info.Participants {
		if p.PUUID == "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg" {
			riotParticipant = p
		}
	}

	require.NotNil(t, riotParticipant)

	mapper := internal.RiotToParticipantMapper{
		Participant:       *riotParticipant,
		MatchID:           riotMatch.Metadata.MatchID,
		MatchDuration:     1131 * time.Second,
		TeamKills:         27,
		TeamGold:          41017,
		TeamDamage:        56169,
		CounterpartGold:   1,
		CounterpartDamage: 1,
	}

	actual := mapper.Map()

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct participant champion",
			Expected: 63,
			Actual:   actual.ChampionID,
		},
		{
			Name:     "expects correct participant kills",
			Expected: 2,
			Actual:   actual.Kills,
		},
		{
			Name:     "expects correct participant kill participation",
			Expected: float32(10.0 / 27.0),
			Actual:   actual.KillParticipation,
		},
		{
			Name:     "expects correct participant cs per minute",
			Expected: float32(131.0 * 60 / 1131),
			Actual:   actual.CreepScorePerMinute,
		},
		{
			Name:     "expects correct participant damage percentage",
			Expected: float32(12629.0 / 56169),
			Actual:   actual.DamagePercentageTeam,
		},
		{
			Name:     "expects correct participant gold percentage",
			Expected: float32(6856.0 / 41017),
			Actual:   actual.GoldPercentageTeam,
		},
		{
			Name:     "expects correct participant position",
			Expected: internal.TeamPositionMiddle,
			Actual:   actual.TeamPosition,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}

func TestRiotToLiveMatchMapper(t *testing.T) {
	riotLiveMatch := sample.WithSampleLiveMatch()

	mapper := internal.RiotToLiveMatchMapper{
		Match: riotLiveMatch,
	}

	actual := mapper.Map()

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
}

func TestRiotToLiveParticipantMapper(t *testing.T) {
	riotLiveMatch := sample.WithSampleLiveMatch()

	riotLiveParticipant := riotLiveMatch.Participants[0]

	mapper := internal.RiotToLiveMatchParticipantMapper{
		Participant: riotLiveParticipant,
		MatchID:     fmt.Sprintf("%s_%d", riotLiveMatch.PlatformID, riotLiveMatch.GameID),
	}

	actual := mapper.Map()

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct participant champion",
			Expected: 517,
			Actual:   actual.ChampionID,
		},
		{
			Name:     "expects correct participant team id",
			Expected: 100,
			Actual:   actual.TeamID,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}

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
