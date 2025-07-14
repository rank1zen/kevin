package postgres_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetChampions(t *testing.T) {
	ctx := context.Background()

	store := DefaultPGInstance.SetupStore(ctx, t)

	p1PUUID := internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

	match := internal.NewMatch(sample.WithSampleMatch())

	match.ID = "M1"
	for i := range match.Participants {
		match.Participants[i].MatchID = "M1"
	}

	match.Date = time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	match.Participants[1].PUUID = p1PUUID
	match.Participants[1].ChampionID = 13
	match.Participants[1].Kills = 2

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
			champions, err := store.GetChampions(ctx, "P1", time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC))
			require.NoError(t, err)

			require.Equal(t, 2, len(champions))
			require.EqualValues(t, 13, champions[0].Champion)
			require.Equal(t, "P1", champions[0].PUUID)

			assert.Equal(t, 2, champions[0].GamesPlayed)
		},
	)

	t.Run(
		"expects kills averaged correctly",
		func(t *testing.T) {
			champions, err := store.GetChampions(ctx, "P1", time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC))
			require.NoError(t, err)

			require.Equal(t, 2, len(champions))
			require.EqualValues(t, 13, champions[0].Champion)
			require.Equal(t, "P1", champions[0].PUUID)

			assert.Equal(t, 2.5, champions[0].Kills)
		},
	)

	t.Run(
		"expects order by games played",
		func(t *testing.T) {
			champions, err := store.GetChampions(ctx, "P1", time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC))
			require.NoError(t, err)

			require.Equal(t, 2, len(champions))

			assert.EqualValues(t, 13, champions[0].Champion)
			assert.EqualValues(t, 12, champions[1].Champion)
		},
	)
}

func TestGetRank(t *testing.T) {
	ctx := context.Background()

	store := DefaultPGInstance.SetupStore(ctx, t)

	puuid := internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

	summoner := internal.Summoner{
		PUUID:      puuid,
		Name:       "T1 OK GOOD YES",
		Tagline:    "NA1",
		Platform:   "NA1",
		SummonerID: "wr1_FUy4RQSAEmUwMiUC8-ttmFopEqhrj4pkyuFXx8ZySs4",
	}

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

	summoner := internal.Summoner{
		PUUID:      puuid,
		Name:       "T1 OK GOOD YES",
		Tagline:    "NA1",
		Platform:   "NA1",
		SummonerID: "wr1_FUy4RQSAEmUwMiUC8-ttmFopEqhrj4pkyuFXx8ZySs4",
	}

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
