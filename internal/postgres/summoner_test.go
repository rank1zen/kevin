package postgres_test

import (
	"context"
	"testing"

	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummonerStore_GetSummoner(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)

	store := postgres.SummonerStore{Tx: pool}

	err := store.CreateSummoner(ctx, postgres.Summoner{PUUID: "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg", Name: "T1 OK GOOD YES", Tagline: "NA1"})
	require.NoError(t, err)

	actual, err := store.GetSummoner(ctx, T1OKGOODYESNA1PUUID)
	require.NoError(t, err)

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct puuid",
			Expected: T1OKGOODYESNA1PUUID,
			Actual:   actual.PUUID,
		},
		{
			Name:     "expects correct name",
			Expected: "T1 OK GOOD YES",
			Actual:   actual.Name,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}

func TestSummonerStore_CreateSummoner(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)

	store := postgres.SummonerStore{Tx: pool}

	err := store.CreateSummoner(ctx, postgres.Summoner{PUUID: "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg", Name: "T1 OK GOOD YES", Tagline: "NA1"})
	if assert.NoError(t, err) {
		_, err := store.GetSummoner(ctx, T1OKGOODYESNA1PUUID)
		assert.NoError(t, err)
	}
}

func TestSummonerStore_SearchByNameTag(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)

	store := postgres.SummonerStore{Tx: pool}

	err := store.CreateSummoner(ctx, postgres.Summoner{
		PUUID:   "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg",
		Name:    "T1 OK GOOD YES",
		Tagline: "NA1",
	})
	require.NoError(t, err)

	t.Run(
		"expects matches prefix",
		func(t *testing.T) {
			actual, err := store.SearchByNameTag(ctx, "T1 ", "")
			if assert.NoError(t, err) {
				assert.Len(t, actual, 1)
			}
		},
	)

	t.Run(
		"expects does not match prefix",
		func(t *testing.T) {
			actual, err := store.SearchByNameTag(ctx, "T1 OK GOOD UE", "")
			if assert.NoError(t, err) {
				assert.Len(t, actual, 0)
			}
		},
	)

	t.Run(
		"expects matches tag",
		func(t *testing.T) {
			actual, err := store.SearchByNameTag(ctx, "", "NA1")
			if assert.NoError(t, err) {
				assert.Len(t, actual, 1)
			}
		},
	)
}
