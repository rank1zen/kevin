package store_test

import (
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParticipantStore_GetParticipantByPUUIDAndMatchIDs(t *testing.T) {
	t.Parallel()

	tx := DefaultPGInstance.SetupTx(t)

	matchStore := store.NewMatchStore(tx)
	participantStore := store.NewParticipantStore(tx)

	match1Date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	match0Date := match1Date.AddDate(0, 0, 1)

	_, err := matchStore.CreateMatch(t.Context(), store.CreateMatch{
		MatchID:  "NA1_123",
		Version:  "26.16.1",
		Date:     match0Date,
		Duration: 20 * time.Minute,
		WinnerID: "100",
	})
	require.NoError(t, err)

	_, err = matchStore.CreateMatch(t.Context(), store.CreateMatch{
		MatchID:  "NA1_456",
		Version:  "26.16.1",
		Date:     match1Date,
		Duration: 21 * time.Minute,
		WinnerID: "200",
	})
	require.NoError(t, err)

	_, err = participantStore.CreateParticipant(t.Context(), store.CreateParticipant{
		MatchID:     "NA1_123",
		PUUID:       ExamplePUUID,
		ItemIDs:     make([]string, 0),
		SummonerIDs: []string{"100", "200"},
		RuneIDs:     make([]string, 0),
	})
	require.NoError(t, err)

	_, err = participantStore.CreateParticipant(t.Context(), store.CreateParticipant{
		MatchID:     "NA1_456",
		PUUID:       ExamplePUUID,
		ItemIDs:     make([]string, 0),
		SummonerIDs: []string{"300", "400"},
		RuneIDs:     make([]string, 0),
	})
	require.NoError(t, err)

	_, err = participantStore.CreateParticipant(t.Context(), store.CreateParticipant{
		MatchID:     "NA1_456",
		PUUID:       ExamplePUUID2,
		ItemIDs:     make([]string, 0),
		SummonerIDs: make([]string, 0),
		RuneIDs:     make([]string, 0),
	})
	require.NoError(t, err)

	t.Run("should return participant details for every match", func(t *testing.T) {
		participants, err := participantStore.GetParticipantByPUUIDAndMatchIDs(t.Context(), ExamplePUUID, []string{"NA1_123", "NA1_456"})
		if assert.NoError(t, err) && assert.Len(t, participants, 2) {
			expectedMatchIDs := []string{"NA1_123", "NA1_456"}
			gotMatchIDs := []string{
				participants[0].MatchID,
				participants[1].MatchID,
			}
			assert.ElementsMatch(t, expectedMatchIDs, gotMatchIDs)

			expectedSummonerSpellIDs := [][]string{{"100", "200"}, {"300", "400"}}
			gotSummonerSpellIDs := [][]string{participants[0].SummonerIDs, participants[1].SummonerIDs}
			assert.ElementsMatch(t, expectedSummonerSpellIDs, gotSummonerSpellIDs)
		}
	})
}
