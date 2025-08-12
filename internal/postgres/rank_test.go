package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRankStore_ListRankIDs(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.RankStore{Tx: pool}

	puuid := internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

	for _, tc := range []postgres.RankStatus{
		{
			RankStatusID:  1,
			PUUID:         puuid.String(),
			EffectiveDate: time.Date(2025, time.April, 0, 0, 0, 0, 0, time.UTC),
			EndDate:       time.Time{},
			IsCurrent:     false,
			IsRanked:      false,
		},
		{
			RankStatusID:  2,
			PUUID:         puuid.String(),
			EffectiveDate: time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC),
			EndDate:       time.Time{},
			IsCurrent:     false,
			IsRanked:      false,
		},
		{
			RankStatusID:  3,
			PUUID:         puuid.String(),
			EffectiveDate: time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC),
			EndDate:       time.Time{},
			IsCurrent:     false,
			IsRanked:      false,
		},
		{
			RankStatusID:  4,
			PUUID:         puuid.String(),
			EffectiveDate: time.Date(2025, time.April, 3, 0, 0, 0, 0, time.UTC),
			EndDate:       time.Time{},
			IsCurrent:     false,
			IsRanked:      false,
		},
	} {
		err := store.CreateRankStatus(ctx, tc)
		require.NoError(t, err)
	}

	for _, tc := range []struct{
		TestName string
		Option postgres.ListRankOption
		Want []int
	}{
		{
			TestName: "returns most recent rank",
			Option:   postgres.ListRankOption{Offset: 0, Limit:  1, Recent: true},
			Want:     []int{4},
		},
	} {
		t.Run(
			tc.TestName,
			func(t *testing.T) {
				ids, err := store.ListRankIDs(ctx, puuid, tc.Option)
				require.NoError(t, err)

				assert.Equal(t, tc.Want, ids)
			},
		)
	}
}
