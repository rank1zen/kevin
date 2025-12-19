package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfileStore_RecordProfile(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)

	store := postgres.NewStore(pool)

	profile := internal.Profile{
		PUUID:   T1OKGOODYESNA1PUUID,
		Name:    "T1 OK GOOD YES",
		Tagline: "NA1",
		Rank:    internal.RankStatus{PUUID: T1OKGOODYESNA1PUUID, EffectiveDate: time.Date(2025, time.April, 4, 0, 0, 0, 0, time.UTC), Detail: nil},
	}

	err := store.Profile.RecordProfile(ctx, &profile)
	assert.NoError(t, err)

	_, err = store.Profile.GetProfile(ctx, T1OKGOODYESNA1PUUID)
	assert.NoError(t, err)
}

func TestProfileStore_GetProfile(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)

	store := postgres.NewStore(pool)

	profile := internal.Profile{
		PUUID:   T1OKGOODYESNA1PUUID,
		Name:    "T1 OK GOOD YES",
		Tagline: "NA1",
		Rank:    internal.RankStatus{PUUID: T1OKGOODYESNA1PUUID, EffectiveDate: time.Date(2025, time.April, 4, 0, 0, 0, 0, time.UTC), Detail: nil},
	}

	err := store.Profile.RecordProfile(ctx, &profile)
	require.NoError(t, err)

	actual, err := store.Profile.GetProfile(ctx, T1OKGOODYESNA1PUUID)
	require.NoError(t, err)

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects unranked",
			Expected: (*internal.RankDetail)(nil),
			Actual:   actual.Rank.Detail,
		},
		{
			Name:     "expects correct puuid",
			Expected: T1OKGOODYESNA1PUUID,
			Actual:   actual.PUUID,
		},
		{
			Name:     "expects correct rank puuid",
			Expected: T1OKGOODYESNA1PUUID,
			Actual:   actual.Rank.PUUID,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}

	t.Run(
		"expects ErrSummonerNotFound",
		func(t *testing.T) {
			_, err := store.Profile.GetProfile(ctx, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q")
			assert.ErrorIs(t, err, internal.ErrSummonerNotFound)
		},
	)
}
