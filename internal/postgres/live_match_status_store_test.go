package postgres_test

import (
	"context"
	"testing"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLiveMatchStatusStore_CreateLiveMatchStatus(t *testing.T) {
	ctx := context.Background()

	store := (*postgres.LiveMatchStatusStore)(DefaultPGInstance.SetupStore(ctx, t))

	err := store.CreateLiveMatchStatus(ctx, &internal.LiveMatchStatus{})
	assert.NoError(t, err)
}

func TestLiveMatchStatusStore_GetLiveMatchStatus(t *testing.T) {
	ctx := context.Background()

	store := (*postgres.LiveMatchStatusStore)(DefaultPGInstance.SetupStore(ctx, t))

	err := store.CreateLiveMatchStatus(ctx, &internal.LiveMatchStatus{})
	require.NoError(t, err)

	status, err := store.GetLiveMatchStatus(ctx, riot.RegionNA1, "NA1")
	assert.Equal(t, nil, status.ID)
	assert.Equal(t, nil, status.Region)
}

func TestLiveMatchStatusStore_ExpireLiveMatch(t *testing.T) {
	ctx := context.Background()

	store := (*postgres.LiveMatchStatusStore)(DefaultPGInstance.SetupStore(ctx, t))

	err := store.CreateLiveMatchStatus(ctx, &internal.LiveMatchStatus{})
	require.NoError(t, err)

	err = store.ExpireLiveMatch(ctx, riot.RegionNA1, "NA1")
	assert.NoError(t, err)

	status, err := store.GetLiveMatchStatus(ctx, riot.RegionNA1, "NA1")
	require.Equal(t, true, status.Expired)
}
