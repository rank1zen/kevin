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

func TestMatchStore_ListMatchHistoryIDs(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.MatchStore{Tx: pool}

	puuid := internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

	for _, tc := range []struct {
		Match       postgres.Match
		Participant postgres.Participant
	}{
		{
			Match: postgres.Match{
				ID:       "1",
				Date:     time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC),
				Duration: 133 * time.Second,
				Version:  "255",
				WinnerID: 100,
			},
			Participant: postgres.Participant{
				MatchID: "1",
				PUUID:   puuid.String(),
				TeamPosition: "Top",
			},
		},
		{
			Match: postgres.Match{
				ID:       "2",
				Date:     time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC),
				Duration: 133 * time.Second,
				Version:  "255",
				WinnerID: 100,
			},
			Participant: postgres.Participant{
				MatchID: "2",
				PUUID:   puuid.String(),
				TeamPosition: "Top",
			},
		},
		{
			Match: postgres.Match{
				ID:       "3",
				Date:     time.Date(2025, time.April, 3, 0, 0, 0, 0, time.UTC),
				Duration: 133 * time.Second,
				Version:  "255",
				WinnerID: 100,
			},
			Participant: postgres.Participant{
				MatchID: "3",
				PUUID:   puuid.String(),
				TeamPosition: "Top",
			},
		},
	} {
		err := store.CreateMatch(ctx, tc.Match)
		require.NoError(t, err)

		err = store.CreateParticipant(ctx, tc.Participant)
		require.NoError(t, err)
	}

	ids, err := store.ListMatchHistoryIDs(ctx, puuid, time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	assert.Equal(t, []string{"1"}, ids)
}

func TestMatchStore_GetMatch(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.MatchStore{Tx: pool}

	expected := postgres.Match{
		ID:       "1",
		Date:     time.Date(2025, 4, 0, 0, 0, 0, 0, time.UTC),
		Duration: 1300 * time.Second,
		Version:  "14.5",
		WinnerID: 100,
	}

	err := store.CreateMatch(ctx, expected)
	require.NoError(t, err)

	actual, err := store.GetMatch(ctx, "1")
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestMatchStore_GetParticipant(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.MatchStore{Tx: pool}

	match := postgres.Match{
		ID:       "1",
		Date:     time.Date(2025, 4, 0, 0, 0, 0, 0, time.UTC),
		Duration: 1300 * time.Second,
		Version:  "14.5",
		WinnerID: 100,
	}

	err := store.CreateMatch(ctx, match)
	require.NoError(t, err)

	puuid := internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

	expected := postgres.Participant{
		PUUID:   puuid.String(),
		MatchID: "1",
		TeamPosition: "Top",
	}

	err = store.CreateParticipant(ctx, expected)
	require.NoError(t, err)

	t.Run(
		"expects equal object",
		func(t *testing.T) {
			actual, err := store.GetParticipant(ctx, puuid, "1")
			require.NoError(t, err)

			assert.Equal(t, expected, actual)
		},
	)
}
