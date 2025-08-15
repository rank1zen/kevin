package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/stretchr/testify/assert"
)

var ExamplePUUID = internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

var (
	ExampleProfileName = "T1 OK GOOD YES"
	ExampleProfileTag  = "NA1"
)

func TestStore_RecordProfile(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)

	store := postgres.NewStore2(pool)

	profile := internal.Profile{
		PUUID:   ExamplePUUID,
		Name:    ExampleProfileName,
		Tagline: ExampleProfileTag,
		Rank:    internal.RankStatus{},
	}

	err := store.RecordProfile(ctx, profile)
	assert.NoError(t, err)

	_, err = store.GetProfileDetail(ctx, ExamplePUUID)
	assert.NoError(t, err)
}

func TestStore_RecordMatch(t *testing.T) {
	ctx := context.Background()

	pool := DefaultPGInstance.SetupConn(ctx, t)

	store := postgres.NewStore2(pool)

	profile := internal.Match{
		ID:           "",
		Date:         time.Time{},
		Duration:     0,
		Version:      "",
		WinnerID:     0,
		Participants: [10]internal.Participant{},
	}

	err := store.RecordMatch(ctx, profile)
	assert.NoError(t, err)

	_, err = store.GetMatchDetail(ctx, "1")
	assert.NoError(t, err)
}
