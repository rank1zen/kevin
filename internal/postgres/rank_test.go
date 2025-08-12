package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRankStore_CreateRankDetail(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.RankStore{Tx: pool}

	puuid := internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

	exampleStatus := postgres.RankStatus{
		PUUID:         puuid.String(),
		EffectiveDate: time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC),
		IsRanked:      false,
	}

	exampleID, err := store.CreateRankStatus(ctx, exampleStatus)
	require.NoError(t, err)

	for _, tc := range []struct {
		TestName string
		Detail   postgres.RankDetail
		Error    error
	}{
		{
			TestName: "expects error for unknown status id",
			Detail: postgres.RankDetail{
				RankStatusID: 9999,
				Wins:         1,
				Losses:       20,
				Tier:         "Iron",
				Division:     "IV",
				LP:           28,
			},
			Error: postgres.ErrInvalidRankStatuID,
		},
		{
			TestName: "expects no error",
			Detail: postgres.RankDetail{
				RankStatusID: exampleID,
				Wins:         1,
				Losses:       20,
				Tier:         "Iron",
				Division:     "IV",
				LP:           28,
			},
			Error: nil,
		},
	} {
		t.Run(
			tc.TestName,
			func(t *testing.T) {
				err := store.CreateRankDetail(ctx, tc.Detail)
				assert.ErrorIs(t, err, tc.Error)
			},
		)
	}
}

func TestRankStore_CreateRankStatus(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.RankStore{Tx: pool}

	puuid := internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

	for _, tc := range []postgres.RankStatus{
		{
			PUUID:         puuid.String(),
			EffectiveDate: time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC),
			IsRanked:      false,
		},
	} {
		_, err := store.CreateRankStatus(ctx, tc)
		require.NoError(t, err)
	}
}

func TestRankStore_GetRankDetail(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.RankStore{Tx: pool}

	puuid := internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

	exampleID, err := store.CreateRankStatus(ctx, postgres.RankStatus{
		PUUID:         puuid.String(),
		EffectiveDate: time.Now(),
		IsRanked:      true,
	})
	require.NoError(t, err)

	exampleDetail := postgres.RankDetail{
		RankStatusID: exampleID,
		Wins:         1,
		Losses:       20,
		Tier:         "Iron",
		Division:     "IV",
		LP:           28,
	}

	err = store.CreateRankDetail(ctx, exampleDetail)
	require.NoError(t, err)

	for _, tc := range []struct {
		TestName string
		StatusID int
		Detail   postgres.RankDetail
		Error    error
	}{
		{
			TestName: "expects no rows error",
			StatusID: 9999,
			Error:    pgx.ErrNoRows,
		},
		{
			TestName: "expects no error",
			Detail:   exampleDetail,
			StatusID: exampleID,
			Error:    nil,
		},
	} {
		t.Run(
			tc.TestName,
			func(t *testing.T) {
				got, err := store.GetRankDetail(ctx, tc.StatusID)
				if assert.ErrorIs(t, err, tc.Error) {
					assert.Equal(t, tc.Detail, got)
				}
			},
		)
	}
}

func TestRankStore_ListRankIDs(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.RankStore{Tx: pool}

	puuid := internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

	ids := []int{}

	for _, tc := range []postgres.RankStatus{
		{
			PUUID:         puuid.String(),
			EffectiveDate: time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC),
			IsRanked:      false,
		},
		{
			PUUID:         puuid.String(),
			EffectiveDate: time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC),
			IsRanked:      false,
		},
		{
			PUUID:         puuid.String(),
			EffectiveDate: time.Date(2025, time.April, 3, 0, 0, 0, 0, time.UTC),
			IsRanked:      false,
		},
		{
			PUUID:         puuid.String(),
			EffectiveDate: time.Date(2025, time.April, 4, 0, 0, 0, 0, time.UTC),
			IsRanked:      false,
		},
	} {
		id, err := store.CreateRankStatus(ctx, tc)
		require.NoError(t, err)
		ids = append(ids, id)
	}

	april3 := time.Date(2025, time.April, 3, 0, 0, 0, 0, time.UTC)

	for _, tc := range []struct {
		TestName string
		Option   postgres.ListRankOption
		Want     []int
	}{
		{
			TestName: "expects most recent rank",
			Option:   postgres.ListRankOption{Offset: 0, Limit: 1, Recent: true},
			Want:     []int{ids[3]},
		},
		{
			TestName: "expects most recent rank before April 3",
			Option: postgres.ListRankOption{
				Offset: 0,
				Limit:  1,
				Recent: true,
				End:    &april3,
			},
			Want: []int{ids[1]},
		},
		{
			TestName: "expects least recent rank after or on April 3",
			Option: postgres.ListRankOption{
				Offset: 0,
				Limit:  1,
				Recent: false,
				Start:  &april3,
			},
			Want: []int{ids[2]},
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
