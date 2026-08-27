package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRankStore_CreateRank(t *testing.T) {
	ctx := context.Background()

	t.Run("should create rank", func(t *testing.T) {
		t.Parallel()

		rankStore := store.NewRankStore(DefaultPGInstance.SetupTx(t))
		lastUpdated := time.Now().Truncate(time.Millisecond)

		rank, err := rankStore.CreateRank(ctx, store.CreateRank{
			PUUID:        ExamplePUUID,
			Wins:         10,
			Losses:       5,
			Tier:         "GOLD",
			Division:     "II",
			LeaguePoints: 42,
			LastUpdated:  lastUpdated,
		})

		if assert.NoError(t, err) {
			assert.Equal(t, ExamplePUUID, rank.PUUID)
			assert.Equal(t, "GOLD", rank.Tier)
			assert.Equal(t, "II", rank.Division)
			assert.EqualValues(t, 42, rank.LeaguePoints)
		}
	})
}

func TestRankStore_GetRankByPUUID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get created rank", func(t *testing.T) {
		t.Parallel()

		rankStore := store.NewRankStore(DefaultPGInstance.SetupTx(t))
		lastUpdated := time.Now().Truncate(time.Millisecond)

		rank, err := rankStore.CreateRank(ctx, store.CreateRank{
			PUUID:        ExamplePUUID,
			Wins:         12,
			Losses:       7,
			Tier:         "PLATINUM",
			Division:     "I",
			LeaguePoints: 55,
			LastUpdated:  lastUpdated,
		})
		require.NoError(t, err)

		got, err := rankStore.GetRankByPUUID(ctx, rank.PUUID)
		if assert.NoError(t, err) {
			assert.Equal(t, rank.PUUID, got.PUUID)
			assert.EqualValues(t, rank.Wins, got.Wins)
			assert.EqualValues(t, lastUpdated, got.LastUpdated.Truncate(time.Millisecond))
		}
	})
}

func TestRankStore_UpdateRankByPUUID(t *testing.T) {
	ctx := context.Background()

	t.Run("should update rank", func(t *testing.T) {
		t.Parallel()

		rankStore := store.NewRankStore(DefaultPGInstance.SetupTx(t))
		lastUpdated := time.Now().Truncate(time.Millisecond)

		rank, err := rankStore.CreateRank(ctx, store.CreateRank{
			PUUID:        ExamplePUUID,
			Wins:         8,
			Losses:       4,
			Tier:         "SILVER",
			Division:     "IV",
			LeaguePoints: 20,
			LastUpdated:  lastUpdated,
		})
		require.NoError(t, err)

		updatedTime := time.Now().Truncate(time.Millisecond)
		updatedRank, err := rankStore.UpdateRankByPUUID(ctx, rank.PUUID, store.UpdateRank{
			Wins:         9,
			Losses:       3,
			Tier:         "GOLD",
			Division:     "III",
			LeaguePoints: 31,
			LastUpdated:  updatedTime,
		})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 9, updatedRank.Wins)
			assert.EqualValues(t, 3, updatedRank.Losses)
			assert.EqualValues(t, updatedTime, updatedRank.LastUpdated.Truncate(time.Millisecond))
		}
	})
}

func TestRankStore_GetRankByPUUIDs(t *testing.T) {
	t.Parallel()

	tx := DefaultPGInstance.SetupTx(t)
	rankStore := store.NewRankStore(tx)

	_, err := rankStore.CreateRank(t.Context(), store.CreateRank{
		PUUID: ExamplePUUID,
	})
	require.NoError(t, err)

	_, err = rankStore.CreateRank(t.Context(), store.CreateRank{
		PUUID: ExamplePUUID2,
	})
	require.NoError(t, err)

	t.Run("should get all ranks", func(t *testing.T) {
		ranks, err := rankStore.GetRankByPUUIDs(t.Context(), []string{ExamplePUUID, ExamplePUUID2})
		if assert.NoError(t, err) && assert.Len(t, ranks, 2) {
			assert.ElementsMatch(t, []string{ranks[0].PUUID, ranks[1].PUUID}, []string{ExamplePUUID, ExamplePUUID2})
		}
	})
}
