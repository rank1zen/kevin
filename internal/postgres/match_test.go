package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/ddragon"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchStore_ListMatchHistoryIDs(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.MatchStore{Tx: pool}

	puuid := riot.PUUID("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

	for _, tc := range []struct {
		Match       postgres.Match
		Participant postgres.Participant
	}{
		{
			Match: postgres.Match{
				ID:       "1",
				Date:     time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC),
				Duration: 133 * time.Second,
				Version:  "255",
				WinnerID: 100,
			},
			Participant: postgres.Participant{
				MatchID:      "1",
				PUUID:        puuid.String(),
				TeamPosition: "Top",
			},
		},
		{
			Match: postgres.Match{
				ID:       "2",
				Date:     time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC),
				Duration: 133 * time.Second,
				Version:  "255",
				WinnerID: 100,
			},
			Participant: postgres.Participant{
				MatchID:      "2",
				PUUID:        puuid.String(),
				TeamPosition: "Top",
			},
		},
		{
			Match: postgres.Match{
				ID:       "3",
				Date:     time.Date(2025, time.April, 3, 0, 0, 0, 0, time.UTC),
				Duration: 133 * time.Second,
				Version:  "255",
				WinnerID: 100,
			},
			Participant: postgres.Participant{
				MatchID:      "3",
				PUUID:        puuid.String(),
				TeamPosition: "Top",
			},
		},
	} {
		err := store.CreateMatch(ctx, tc.Match)
		require.NoError(t, err)

		err = store.CreateParticipant(ctx, tc.Participant)
		require.NoError(t, err)
	}

	for _, tc := range []struct {
		Name string

		Start, End time.Time

		Expected []string
	}{
		{
			Name:     "expects exclusive end time",
			Start:    time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC),
			End:      time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC),
			Expected: []string{"1"},
		},
		{
			Name:     "expects timezone conversion",
			Start:    time.Date(2025, time.March, 31, 20, 0, 0, 0, Toronto),
			End:      time.Date(2025, time.April, 1, 20, 0, 0, 0, Toronto),
			Expected: []string{"1"},
		},
	} {
		t.Run(
			tc.Name,
			func(t *testing.T) {
				ids, err := store.ListMatchHistoryIDs(ctx, puuid, tc.Start, tc.End)
				require.NoError(t, err)
				assert.Equal(t, tc.Expected, ids)
			},
		)
	}
}

func TestMatchStore_GetMatch(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.MatchStore{Tx: pool}

	expected := postgres.Match{
		ID:       "1",
		Date:     time.Date(2025, 4, 0, 0, 0, 0, 0, time.UTC),
		Duration: 1300 * time.Second,
		Version:  "14.5",
		WinnerID: 100,
	}

	err := store.CreateMatch(ctx, expected)
	require.NoError(t, err)

	actual, err := store.GetMatch(ctx, "1")
	require.NoError(t, err)

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct id",
			Expected: "1",
			Actual:   actual.ID,
		},
		{
			Name:     "expects correct version",
			Expected: "14.5",
			Actual:   actual.Version,
		},
		{
			Name:     "expects correct date",
			Expected: time.Date(2025, 4, 0, 0, 0, 0, 0, time.UTC),
			Actual:   actual.Date.In(time.UTC),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}

func TestMatchStore_GetParticipant(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.MatchStore{Tx: pool}

	riotMatch := sample.Match5346312088()

	match := postgres.Match{
		ID:       riotMatch.Metadata.MatchID,
		Date:     time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		Duration: 6001 * time.Second,
		Version:  riotMatch.Info.GameVersion,
		WinnerID: 100,
	}

	err := store.CreateMatch(ctx, match)
	require.NoError(t, err)

	expected := postgres.Participant{
		PUUID:                T1OKGOODYESNA1PUUID.String(),
		MatchID:              "NA1_5346312088",
		TeamPosition:         "Top",
		DamagePercentageTeam: 1.0 / 3,
	}

	err = store.CreateParticipant(ctx, expected)
	require.NoError(t, err)

	actual, err := store.GetParticipant(ctx, T1OKGOODYESNA1PUUID, "NA1_5346312088")
	require.NoError(t, err)

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct participant damage percentage team",
			Expected: float32(1.0 / 3),
			Actual:   actual.DamagePercentageTeam,
		},
		{
			Name:     "expects correct participant puuid",
			Expected: T1OKGOODYESNA1PUUID.String(),
			Actual:   actual.PUUID,
		},
		{
			Name:     "expects correct participant position",
			Expected: "Top",
			Actual:   actual.TeamPosition,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}

func TestMatchStore_CreateMatch(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.MatchStore{Tx: pool}

	riotMatch := sample.Match5346312088()

	match := postgres.Match{
		ID:       riotMatch.Metadata.MatchID,
		Date:     time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		Duration: 6001 * time.Second,
		Version:  riotMatch.Info.GameVersion,
		WinnerID: 100,
	}

	err := store.CreateMatch(ctx, match)
	if assert.NoError(t, err) {
		_, err := store.GetMatch(ctx, "NA1_5346312088")
		assert.NoError(t, err)
	}
}

func TestMatchStore_CreateParticipant(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.MatchStore{Tx: pool}

	riotMatch := sample.Match5346312088()

	match := postgres.Match{
		ID:       riotMatch.Metadata.MatchID,
		Date:     time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		Duration: 6001 * time.Second,
		Version:  riotMatch.Info.GameVersion,
		WinnerID: 100,
	}

	err := store.CreateMatch(ctx, match)
	require.NoError(t, err)

	err = store.CreateParticipant(ctx, postgres.Participant{
		PUUID:                T1OKGOODYESNA1PUUID.String(),
		MatchID:              "NA1_5346312088",
		TeamID:               100,
		ChampionID:           33,
		ChampionLevel:        18,
		TeamPosition:         "Top",
		SummonerIDs:          [2]int{},
		Runes:                [11]int{},
		Items:                [7]int{},
		Kills:                1,
		Deaths:               2,
		Assists:              3,
		KillParticipation:    4,
		CreepScore:           5,
		CreepScorePerMinute:  6.7,
		DamageDealt:          7,
		DamageTaken:          8,
		DamageDeltaEnemy:     9,
		DamagePercentageTeam: 0.67,
		GoldEarned:           11,
		GoldDeltaEnemy:       12,
		GoldPercentageTeam:   0.41,
		VisionScore:          14,
		PinkWardsBought:      15,
	})

	if assert.NoError(t, err) {
		_, err := store.GetParticipant(ctx, T1OKGOODYESNA1PUUID, "NA1_5346312088")
		assert.NoError(t, err)
	}
}

func TestMatchStore_GetSummonerChampions(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)

	store := postgres.NewStore(pool)

	for _, m := range []riot.Match{
		sample.Match5347748140(),
		sample.Match5347728946(),
		sample.Match5346312088(),
	} {
		mapper := internal.RiotToMatchMapper{
			Match: m,
		}

		match := mapper.Map()
		err := store.Match.RecordMatch(ctx, match)
		require.NoError(t, err)
	}

	matchStore := postgres.MatchStore{Tx: pool}

	actual, err := matchStore.GetSummonerChampions(ctx, T1OKGOODYESNA1PUUID, time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	order := []internal.Champion{}
	for _, champion := range actual {
		order = append(order, champion.Champion)
	}

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects 2 champions played",
			Expected: 2,
			Actual:   len(actual),
		},
		{
			Name:     "expects correct order",
			Expected: []internal.Champion{ddragon.ChampionUrgotID, ddragon.ChampionIllaoiID},
			Actual:   order,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}

	require.Len(t, actual, 2)
	urgot := actual[0]

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects 2 games played",
			Expected: 2,
			Actual:   urgot.GamesPlayed,
		},
		{
			Name:     "expects 1 win",
			Expected: 1,
			Actual:   urgot.Wins,
		},
		{
			Name:     "expects 9.5 kills per game",
			Expected: float32(19) / 2,
			Actual:   urgot.AverageKillsPerGame,
		},
		{
			Name:     "expects 0 pink wards bought per game",
			Expected: float32(0),
			Actual:   urgot.AveragePinkWardsBoughtPerGame,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}
