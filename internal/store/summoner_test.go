package store_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal/pgtestcontainer"
	"github.com/rank1zen/kevin/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ExamplePUUID = "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg"

var DefaultPGInstance *pgtestcontainer.PGInstance

func TestMain(m *testing.M) {
	DefaultPGInstance = pgtestcontainer.NewPGInstance(context.Background())
	defer func(DefaultPGInstance *pgtestcontainer.PGInstance, ctx context.Context) {
		err := DefaultPGInstance.Terminate(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(DefaultPGInstance, context.Background())
	m.Run()
}

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
			Name:          "test-name",
			Tagline:       "test-tagline",
			SummonerLevel: 1,
			ProfileIconID: "test-profile-icon-id",
		})

		require.NoError(t, err)

		updateTime := time.Now().Truncate(time.Millisecond)
		updatedSummoner, err := summonerStore.UpdateSummonerByPUUID(ctx, summoner.PUUID, store.UpdateSummoner{
			Name:          "new-name",
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
