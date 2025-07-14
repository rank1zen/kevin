package internal_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/ddragon"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMatch(t *testing.T) {
	testdata := os.DirFS("../testdata")

	matchFile, err := testdata.Open("NA1_5304757838.json")
	require.NoError(t, err)

	var riotMatch riot.Match
	err = json.NewDecoder(matchFile).Decode(&riotMatch)
	require.NoError(t, err)

	t.Run(
		"create match from riot",
		func(t *testing.T) {
			expected := internal.Match{
				ID:       "NA1_5304757838",
				Date:     time.UnixMilli(1749596377340),
				Duration: 1131 * time.Second,
				Version:  "15.11.685.5259",
				WinnerID: 100,
			}
			actual := internal.NewMatch(internal.WithRiotMatch(&riotMatch))
			assert.Equal(t, expected, actual)
		},
	)
}

func TestParticipant(t *testing.T) {
	testdata := os.DirFS("../testdata")

	matchFile, err := testdata.Open("NA1_5304757838.json")
	require.NoError(t, err)
	var riotMatch riot.Match
	err = json.NewDecoder(matchFile).Decode(&riotMatch)
	require.NoError(t, err)

	expected := internal.Participant{
		PUUID:         internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg"),
		MatchID:       "NA1_5304757838",
		TeamID:        100,
		ChampionID:    63,
		ChampionLevel: 12,
		SummonerIDs:   [2]int{ddragon.SummonerFlashID, ddragon.SummonerDotID},
		Runes: internal.RunePage{
			PrimaryTree:     ddragon.RuneTreeSorceryID,
			PrimaryKeystone: ddragon.RuneArcaneComet.ID,
			PrimaryA:        ddragon.RuneManaflowBand.ID,
			PrimaryB:        ddragon.RuneTranscendence.ID,
			PrimaryC:        ddragon.RuneGatheringStorm.ID,
			SecondaryTree:   ddragon.RuneTreePrecisionID,
			SecondaryA:      ddragon.RunePresenceOfMind.ID,
			SecondaryB:      ddragon.RuneCoupDeGrace.ID,
			MiniOffense:     5005,
			MiniFlex:        5008,
			MiniDefense:     5001,
		},
		Items:                [7]int{1056, 3116, 3020, 2508, 3802, 0, 3363},
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

	t.Run(
		"create participant from riot",
		func(t *testing.T) {
			actual := internal.NewParticipant(internal.RiotMatchToParticipant(riotMatch, "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg"))
			assert.Equal(t, expected, actual)
		},
	)
}

func TestLiveParticipant(t *testing.T) {
	testdata := os.DirFS("../testdata")

	matchFile, err := testdata.Open("spectator/aram.json")

	var riotMatch riot.LiveMatch
	err = json.NewDecoder(matchFile).Decode(&riotMatch)
	require.NoError(t, err)

	t.Run(
		"create live participant from riot",
		func(t *testing.T) {
			expected := internal.LiveParticipant{
				PUUID:        "FJqkUwySDuAPRCwxS2e_WTHYJ9rCXGTE-usG_Rya-rzxDSplTW3FQ2oPp0kev2FlxKV26A3O917gvg",
				MatchID:      "NA1_5308207011",
				ChampionID:   107,
				Runes:        internal.NewRunePage(internal.WithIntList([11]int{8000,8010,9111,9105,8299,8100,8143,8135,5005,5008,5001})),
				TeamID:       100,
				SummonersIDs: [2]int{4, 14},
			}

			actual := internal.NewLiveParticipant(internal.WithRiotCurrentGame(riotMatch, "FJqkUwySDuAPRCwxS2e_WTHYJ9rCXGTE-usG_Rya-rzxDSplTW3FQ2oPp0kev2FlxKV26A3O917gvg"))

			assert.Equal(t, expected, actual)
		},
	)
}
