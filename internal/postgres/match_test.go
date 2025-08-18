package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchStore_ListMatchHistoryIDs(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.MatchStore{Tx: pool}

	puuid := internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

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

	ids, err := store.ListMatchHistoryIDs(ctx, puuid, time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	assert.Equal(t, []string{"1"}, ids)
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

	assert.Equal(t, expected, actual)
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
