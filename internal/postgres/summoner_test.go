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

	err := store.CreateSummoner(ctx, postgres.Summoner{PUUID: T1OKGOODYESNA1PUUID, Name: "T1 OK GOOD YES", Tagline: "NA1"})
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

	err := store.CreateSummoner(ctx, postgres.Summoner{PUUID: T1OKGOODYESNA1PUUID, Name: "T1 OK GOOD YES", Tagline: "NA1"})
	if assert.NoError(t, err) {
		_, err := store.GetSummoner(ctx, T1OKGOODYESNA1PUUID)
		assert.NoError(t, err)
	}
}
