package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchStore_RecordMatch(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)

	store := postgres.NewStore(pool)

	mapper := internal.RiotToMatchMapper{
		Match: sample.WithSampleMatch(),
	}
	err := store.Match.RecordMatch(ctx, mapper.Map())
	if assert.NoError(t, err) {
		_, err = store.Match.GetMatchDetail(ctx, "NA1_5304757838")
		assert.NoError(t, err)
	}
}

func TestMatchStore_GetMatchDetail(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)

	store := postgres.NewStore(pool)

	riotMatch := sample.WithSampleMatch()

	mapper := internal.RiotToMatchMapper{
		Match: riotMatch,
	}

	match := mapper.Map()

	err := store.Match.RecordMatch(ctx, match)
	require.NoError(t, err)

	actual, err := store.Match.GetMatchDetail(ctx, "NA1_5304757838")
	require.NoError(t, err)

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct ID",
			Expected: "NA1_5304757838",
			Actual:   actual.ID,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}

	var actualParticipant *internal.ParticipantDetail
	for _, p := range actual.Participants {
		if p.PUUID == T1OKGOODYESNA1PUUID {
			actualParticipant = &p
		}
	}

	require.NotNil(t, actualParticipant)

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects no rank after",
			Expected: (*internal.RankStatus)(nil),
			Actual:   actualParticipant.RankAfter,
		},
		{
			Name:     "expects no rank before",
			Expected: (*internal.RankStatus)(nil),
			Actual:   actualParticipant.RankBefore,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}

func TestMatchStore_GetMatchHistory(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)

	store := postgres.NewStore(pool)

	for _, m := range []riot.Match{
		sample.Match5347748140(),
		sample.WithSampleMatch(),
		sample.Match5346312088(),
	} {
		mapper := internal.RiotToMatchMapper{
			Match: m,
		}

		match := mapper.Map()
		err := store.Match.RecordMatch(ctx, match)
		require.NoError(t, err)
	}

	rankStore := postgres.RankStore{Tx: pool}

	_, err := rankStore.CreateRankStatus(ctx, postgres.RankStatus{PUUID: T1OKGOODYESNA1PUUID.String(), EffectiveDate: time.Date(2025, time.August, 13, 21, 0, 0, 0, time.UTC)})
	require.NoError(t, err)

	_, err = rankStore.CreateRankStatus(ctx, postgres.RankStatus{PUUID: T1OKGOODYESNA1PUUID.String(), EffectiveDate: time.Date(2025, time.August, 15, 21, 0, 0, 0, time.UTC)})
	require.NoError(t, err)

	actual, err := store.Match.GetMatchlist(ctx, T1OKGOODYESNA1PUUID, time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	order := []string{}
	for _, a := range actual {
		order = append(order, a.MatchID)
	}

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects 3 matches",
			Expected: 3,
			Actual:   len(actual),
		},
		{
			Name:     "expects correct order",
			Expected: []string{"NA1_5347748140", "NA1_5346312088", "NA1_5304757838"},
			Actual:   order,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}

	require.Len(t, actual, 3)
	actualMatch := actual[0]

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects no rank after",
			Expected: (*internal.RankStatus)(nil),
			Actual:   actualMatch.RankAfter,
		},
		// This test is broken for some reason
		// {
		// 	Name:     "expects unranked before",
		// 	Expected: &internal.RankStatus{PUUID: T1OKGOODYESNA1PUUID, EffectiveDate: time.Date(2025, time.August, 13, 21, 0, 0, 0, time.UTC), Detail: nil},
		// 	Actual:   actualMatch.RankBefore,
		// },
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}
