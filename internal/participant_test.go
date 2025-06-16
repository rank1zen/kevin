package internal_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/ddragon"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	match5304757838 = internal.Participant{
		Puuid:         "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg",
		MatchID:       "NA1_5304757838",
		Team:          100,
		Champion:      63,
		ChampionLevel: 12,
		Summoners: [2]internal.Spell{
			internal.Spell(ddragon.SummonerFlash.ID),
			internal.Spell(ddragon.SummonerDot.ID),
		},
		Runes: internal.RunePage{
			PrimaryTree:     internal.Rune(ddragon.Sorcery.ID),
			PrimaryKeystone: internal.Rune(ddragon.ArcaneComet.ID),
			PrimaryA:        internal.Rune(ddragon.ManaflowBand.ID),
			PrimaryB:        internal.Rune(ddragon.Transcendence.ID),
			PrimaryC:        internal.Rune(ddragon.GatheringStorm.ID),
			SecondaryTree:   internal.Rune(ddragon.Precision.ID),
			SecondaryA:      internal.Rune(ddragon.PresenceOfMind.ID),
			SecondaryB:      internal.Rune(ddragon.CoupDeGrace.ID),
			MiniOffense:     5005,
			MiniFlex:        5008,
			MiniDefense:     5001,
		},
		Items: [7]internal.Item{
			internal.Item(1056),
			internal.Item(3116),
			internal.Item(3020),
			internal.Item(2508),
			internal.Item(3802),
			internal.Item(0),
			internal.Item(3363),
		},
		Kills:                2,
		Deaths:               0,
		Assists:              8,
		KillParticipation:    10.0 / 27.0,
		CreepScore:           131,
		CreepScorePerMinute:  131.0 * 60 / 1131,
		DamageDealt:          12629,
		DamageTaken:          10465,
		DamageDeltaEnemy:     7095,
		DamagePercentageTeam: 12629.0 / 56169,
		GoldEarned:           6856,
		GoldDeltaEnemy:       715,
		GoldPercentageTeam:   6856.0 / 41017,
		VisionScore:          7,
		PinkWardsBought:      0,
	}
)

func TestParticipant(t *testing.T) {
	testdata := os.DirFS("../testdata")

	matchFile, err := testdata.Open("NA1_5304757838.json")
	require.NoError(t, err)
	var riotMatch riot.Match
	err = json.NewDecoder(matchFile).Decode(&riotMatch)
	require.NoError(t, err)

	expected := match5304757838

	t.Run(
		"create participant from riot",
		func(t *testing.T) {
			actual := internal.NewParticipant(internal.RiotMatchToParticipant(riotMatch, "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg"))
			assert.Equal(t, expected, actual)
		},
	)
}
