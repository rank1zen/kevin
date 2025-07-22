package postgres_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetChampions(t *testing.T) {
	ctx := context.Background()

	store := DefaultPGInstance.SetupStore(ctx, t)

	p1PUUID := internal.NewPUUIDFromString("0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q")

	match := internal.NewMatch(sample.WithSampleMatch())

	match.ID = "M1"
	for i := range match.Participants {
		match.Participants[i].MatchID = "M1"
	}

	match.Date = time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	match.Participants[1].PUUID = p1PUUID
	match.Participants[1].ChampionID = 13
	match.Participants[1].Kills = 2
	match.WinnerID = match.Participants[1].TeamID

	err := store.RecordMatch(ctx, match)
	require.NoError(t, err)

	match = internal.NewMatch(sample.WithSampleMatch())

	match.ID = "M2"
	for i := range match.Participants {
		match.Participants[i].MatchID = "M2"
	}

	match.Date = time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC)
	match.Participants[1].PUUID = p1PUUID
	match.Participants[1].ChampionID = 13
	match.Participants[1].Kills = 3
	if match.Participants[1].TeamID == 100 {
		match.WinnerID = 200
	} else {
		match.WinnerID = 100
	}

	err = store.RecordMatch(ctx, match)
	require.NoError(t, err)

	match = internal.NewMatch(sample.WithSampleMatch())

	match.ID = "M3"
	for i := range match.Participants {
		match.Participants[i].MatchID = "M3"
	}

	match.Date = time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC)
	match.Participants[1].PUUID = p1PUUID
	match.Participants[1].ChampionID = 12

	err = store.RecordMatch(ctx, match)
	require.NoError(t, err)

	t.Run(
		"expects inclusive date range",
		func(t *testing.T) {
			champions, err := store.GetChampions(ctx, p1PUUID, time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC))
			require.NoError(t, err)

			require.Equal(t, 2, len(champions))
			require.EqualValues(t, 13, champions[0].Champion)
			require.Equal(t, p1PUUID, champions[0].PUUID)

			assert.Equal(t, 2, champions[0].GamesPlayed)
		},
	)

	t.Run(
		"expects order by games played",
		func(t *testing.T) {
			champions, err := store.GetChampions(ctx, p1PUUID, time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC))
			require.NoError(t, err)

			require.Equal(t, 2, len(champions))

			assert.EqualValues(t, 13, champions[0].Champion)
			assert.EqualValues(t, 12, champions[1].Champion)
		},
	)

	champions, err := store.GetChampions(ctx, p1PUUID, time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Equal(t, 2, len(champions))
	require.EqualValues(t, 13, champions[0].Champion)
	ryze := champions[0]

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct number of games for Ryze",
			Expected: 2,
			Actual:   ryze.GamesPlayed,
		},
		{
			Name:     "expects correct number of wins for Ryze",
			Expected: 1,
			Actual:   ryze.Wins,
		},
		{
			Name:     "expects correct kill average for Ryze",
			Expected: 2.5,
			Actual:   ryze.AverageKillsPerGame,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.EqualValues(t, tc.Expected, tc.Actual) })
	}
}

func TestGetRank(t *testing.T) {
	ctx := context.Background()

	store := DefaultPGInstance.SetupStore(ctx, t)

	puuid := internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

	summoner := sample.SummonerT1OKGOODYESNA1

	for _, rank := range []internal.RankStatus{
		{
			PUUID:         puuid,
			EffectiveDate: time.Date(2025, time.April, 0, 0, 0, 0, 0, time.UTC),
			Detail:        nil,
		},
		{
			PUUID:         puuid,
			EffectiveDate: time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC),
			Detail:        nil,
		},
		{
			PUUID:         puuid,
			EffectiveDate: time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC),
			Detail:        nil,
		},
		{
			PUUID:         puuid,
			EffectiveDate: time.Date(2025, time.April, 3, 0, 0, 0, 0, time.UTC),
			Detail:        nil,
		},
	} {
		err := store.RecordSummoner(ctx, summoner, rank)
		require.NoError(t, err)
	}

	t.Run(
		"returns most recent rank",
		func(t *testing.T) {
			rank, err := store.GetRank(ctx, puuid, time.Now(), true)
			require.NoError(t, err)

			var endDate *time.Time = nil

			assert.Equal(t, time.Date(2025, time.April, 3, 0, 0, 0, 0, time.UTC), rank.EffectiveDate)
			assert.Equal(t, endDate, rank.EndDate)
			assert.Equal(t, true, rank.IsCurrent)
		},
	)

	t.Run(
		"returns rank record at date",
		func(t *testing.T) {
			rank, err := store.GetRank(ctx, puuid, time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC), true)
			require.NoError(t, err)

			endDate := time.Date(2025, time.April, 3, 0, 0, 0, 0, time.UTC)

			assert.Equal(t, time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC), rank.EffectiveDate)
			assert.Equal(t, &endDate, rank.EndDate)
			assert.Equal(t, false, rank.IsCurrent)
		},
	)

	t.Run(
		"returns rank record after date",
		func(t *testing.T) {
			rank, err := store.GetRank(ctx, puuid, time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC), false)
			require.NoError(t, err)

			assert.Equal(t, time.Date(2025, time.April, 3, 0, 0, 0, 0, time.UTC), rank.EffectiveDate)
		},
	)
}

func TestGetMatches(t *testing.T) {
	// ctx := context.Background()
	//
	// store := DefaultPGInstance.SetupStore(ctx, t)
	//
	// puuid := ""
	// matchIDs := []string{}
	//
	// for _, id := range matchIDs {
	// 	err := store.RecordMatch(ctx, nil, nil)
	// }
	//
	// t.Run(
	// 	"returns in chronological order",
	// 	func(t *testing.T) {
	// 		matches, err := store.GetMatches(ctx, puuid, 0)
	// 		require.NoError(t, err)
	//
	// 		actual := []string{}
	// 		for _, match := range matches {
	// 			actual = append(actual, match.MatchID)
	// 		}
	//
	// 		expected := []string{}
	//
	// 		assert.Equal(t, expected, actual)
	// 	},
	// )
}

func TestGetSummoner(t *testing.T) {
	ctx := context.Background()

	store := DefaultPGInstance.SetupStore(ctx, t)

	puuid := internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

	t.Run(
		"fails when summoner is not found",
		func(t *testing.T) {
			_, err := store.GetSummoner(ctx, puuid)
			assert.ErrorIs(t, err, internal.ErrSummonerNotFound)
		},
	)

	summoner := sample.SummonerT1OKGOODYESNA1

	rank := internal.RankStatus{
		PUUID:         puuid,
		EffectiveDate: time.Now(),
		Detail:        nil,
	}

	err := store.RecordSummoner(ctx, summoner, rank)
	require.NoError(t, err)

	t.Run(
		"returns correct values",
		func(t *testing.T) {
			actual, err := store.GetSummoner(ctx, puuid)
			require.NoError(t, err)
			assert.Equal(t, summoner, actual)
		},
	)
}

var DefaultPGInstance *postgres.PGInstance

func TestMain(t *testing.M) {
	ctx := context.Background()

	DefaultPGInstance = postgres.NewPGInstance(context.Background())

	code := t.Run()

	if err := DefaultPGInstance.Terminate(ctx); err != nil {
		log.Fatalf("terminating: %s", err)
	}

	os.Exit(code)
}
