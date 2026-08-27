package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummonerStore_CreateSummoner(t *testing.T) {
	ctx := context.Background()

	t.Run("should create summoner", func(t *testing.T) {
		t.Parallel()

		summonerStore := store.NewSummonerStore(DefaultPGInstance.SetupTx(t))

		summoner, err := summonerStore.CreateSummoner(ctx, store.CreateSummoner{
			PUUID:         ExamplePUUID,
			Name:          "test-name",
			Tagline:       "test-tagline",
			SummonerLevel: 1,
			ProfileIconID: "test-profile-icon-id",
			LastUpdated:   time.Now().Truncate(time.Millisecond),
		})

		if assert.NoError(t, err) {
			assert.Equal(t, ExamplePUUID, summoner.PUUID)
			assert.Equal(t, "test-profile-icon-id", summoner.ProfileIconID)
		}
	})
}

func TestSummonerStore_GetSummonerByPUUID(t *testing.T) {
	ctx := context.Background()

	t.Run("should get created summoner", func(t *testing.T) {
		t.Parallel()

		summonerStore := store.NewSummonerStore(DefaultPGInstance.SetupTx(t))

		summoner, err := summonerStore.CreateSummoner(ctx, store.CreateSummoner{
			PUUID:         ExamplePUUID,
			Name:          "test-name",
			Tagline:       "test-tagline",
			SummonerLevel: 1,
			ProfileIconID: "test-profile-icon-id",
		})

		require.NoError(t, err)

		got, err := summonerStore.GetSummonerByPUUID(ctx, summoner.PUUID)
		if assert.NoError(t, err) {
			assert.Equal(t, summoner.PUUID, got.PUUID)
			assert.Equal(t, summoner.ProfileIconID, got.ProfileIconID)
		}
	})
}

func TestSummonerStore_UpdateSummonerByPUUID(t *testing.T) {
	ctx := context.Background()

	t.Run("should update summoner", func(t *testing.T) {
		t.Parallel()

		summonerStore := store.NewSummonerStore(DefaultPGInstance.SetupTx(t))

		summoner, err := summonerStore.CreateSummoner(ctx, store.CreateSummoner{
			PUUID:         ExamplePUUID,
			Region:        "test-region",
			Name:          "test-name",
			Tagline:       "test-tagline",
			SummonerLevel: 1,
			ProfileIconID: "test-profile-icon-id",
		})

		require.NoError(t, err)

		updateTime := time.Now().Truncate(time.Millisecond)
		updatedSummoner, err := summonerStore.UpdateSummonerByPUUID(ctx, summoner.PUUID, store.UpdateSummoner{
			Name:          "new-name",
			Region:        summoner.Region,
			Tagline:       summoner.Tagline,
			SummonerLevel: summoner.SummonerLevel,
			ProfileIconID: summoner.ProfileIconID,
			LastUpdated:   updateTime,
		})
		if assert.NoError(t, err) {
			assert.EqualValues(t, "new-name", updatedSummoner.Name)
			assert.EqualValues(t, summoner.Tagline, updatedSummoner.Tagline)
			assert.EqualValues(t, updateTime, updatedSummoner.LastUpdated.Truncate(time.Millisecond))
		}
	})
}

func TestSummonerStore_GetSummonerByRegionNameTag(t *testing.T) {
	t.Parallel()

	tx := DefaultPGInstance.SetupTx(t)

	t.Run("should get summoner by region, name, and tag", func(t *testing.T) {
		summonerStore := store.NewSummonerStore(tx)

		_, err := summonerStore.CreateSummoner(t.Context(), store.CreateSummoner{
			PUUID:         ExamplePUUID,
			Region:        "test-region",
			Name:          "test-name",
			Tagline:       "test-tagline",
			SummonerLevel: 1,
			ProfileIconID: "test-profile-icon-id",
			LastUpdated:   time.Now().Truncate(time.Millisecond),
		})
		require.NoError(t, err)

		got, err := summonerStore.GetSummonerByRegionNameTag(t.Context(), "test-region", "test-name", "test-tagline")
		if assert.NoError(t, err) {
			assert.Equal(t, ExamplePUUID, got.PUUID)
		}
	})
}
