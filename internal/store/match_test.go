package store_test

import (
	"context"
	"testing"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/kevin/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchStore_GetMatchByPUUID(t *testing.T) {
	ctx := context.Background()

	tx := DefaultPGInstance.SetupTx(t)

	nestedTx := func(tb testing.TB) pgx.Tx {
		tx, err := tx.Begin(tb.Context())
		if err != nil {
			t.Fatal(err)
		}

		tb.Cleanup(func() {
			err := tx.Rollback(context.Background())
			if err != nil {
				t.Fatal(err)
			}
		})

		return tx
	}

	matchStore := store.NewMatchStore(tx)
	participantStore := store.NewParticipantStore(tx)

	match1Date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	match0Date := match1Date.AddDate(0, 0, 1)

	_, err := matchStore.CreateMatch(ctx, store.CreateMatch{
		MatchID:  "NA1_123",
		Version:  "26.16.1",
		Date:     match0Date,
		Duration: 20 * time.Minute,
		WinnerID: "100",
	})
	require.NoError(t, err)

	_, err = matchStore.CreateMatch(ctx, store.CreateMatch{
		MatchID:  "NA1_456",
		Version:  "26.16.1",
		Date:     match1Date,
		Duration: 21 * time.Minute,
		WinnerID: "200",
	})
	require.NoError(t, err)

	_, err = participantStore.CreateParticipant(ctx, store.CreateParticipant{
		MatchID:     "NA1_123",
		PUUID:       ExamplePUUID,
		ItemIDs:     make([]string, 0),
		SummonerIDs: make([]string, 0),
		RuneIDs:     make([]string, 0),
	})
	require.NoError(t, err)

	_, err = participantStore.CreateParticipant(ctx, store.CreateParticipant{
		MatchID:     "NA1_456",
		PUUID:       ExamplePUUID,
		ItemIDs:     make([]string, 0),
		SummonerIDs: make([]string, 0),
		RuneIDs:     make([]string, 0),
	})
	require.NoError(t, err)

	_, err = participantStore.CreateParticipant(ctx, store.CreateParticipant{
		MatchID:     "NA1_456",
		PUUID:       ExamplePUUID2,
		ItemIDs:     make([]string, 0),
		SummonerIDs: make([]string, 0),
		RuneIDs:     make([]string, 0),
	})
	require.NoError(t, err)

	t.Run("should find all matches", func(t *testing.T) {
		tx := nestedTx(t)

		matchStore := store.NewMatchStore(tx)

		matches, nextPage, err := matchStore.GetMatchByPUUID(ctx, ExamplePUUID, store.GetMatchByPUUIDPageParams{})
		if assert.NoError(t, err) {
			assert.Equal(t, 2, len(matches))
			assert.EqualValues(t, uuid.Nil(), nextPage.LastID)
			assert.EqualValues(t, time.Time{}, nextPage.LastDate)
		}
	})

	t.Run("should page correctly", func(t *testing.T) {
		tx := nestedTx(t)

		matchStore := store.NewMatchStore(tx)

		matches, nextPage, err := matchStore.GetMatchByPUUID(ctx, ExamplePUUID, store.GetMatchByPUUIDPageParams{
			Size: 1,
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(matches))
		require.EqualValues(t, match0Date, nextPage.LastDate.In(time.UTC))

		matches, nextPage, err = matchStore.GetMatchByPUUID(ctx, ExamplePUUID, nextPage)
		require.NoError(t, err)
		require.Equal(t, 1, len(matches))
		require.EqualValues(t, match1Date, nextPage.LastDate.In(time.UTC))
	})
}
