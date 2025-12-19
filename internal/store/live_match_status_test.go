package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLiveMatchStatusStore_CreateLiveMatchStatus(t *testing.T) {
	ctx := context.Background()

	store := (*store.LiveMatchStatusStore)(setupStore(ctx, t))

	err := store.CreateLiveMatchStatus(ctx, &internal.LiveMatchStatus{
		Region:  riot.RegionNA1,
		ID:      "NA1_1234567890",
		Date:    time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		Expired: false,
	})
	assert.NoError(t, err)
}

func TestLiveMatchStatusStore_GetLiveMatchStatus(t *testing.T) {
	ctx := context.Background()

	store := (*store.LiveMatchStatusStore)(setupStore(ctx, t))

	err := store.CreateLiveMatchStatus(ctx, &internal.LiveMatchStatus{
		Region:  riot.RegionNA1,
		ID:      "NA1_1234567890",
		Date:    time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		Expired: false,
	})
	require.NoError(t, err)

	status, err := store.GetLiveMatchStatus(ctx, riot.RegionNA1, "NA1_1234567890")
	if assert.NoError(t, err) {
		assert.Equal(t, "NA1_1234567890", status.ID)
		assert.Equal(t, riot.RegionNA1, status.Region)
	}
}

func TestLiveMatchStatusStore_ExpireLiveMatch(t *testing.T) {
	ctx := context.Background()

	store := (*store.LiveMatchStatusStore)(setupStore(ctx, t))

	err := store.CreateLiveMatchStatus(ctx, &internal.LiveMatchStatus{
		Region:  riot.RegionNA1,
		ID:      "NA1_1234567890",
		Date:    time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		Expired: false,
	})
	require.NoError(t, err)

	err = store.ExpireLiveMatch(ctx, riot.RegionNA1, "NA1_1234567890")
	assert.NoError(t, err)

	status, err := store.GetLiveMatchStatus(ctx, riot.RegionNA1, "NA1_1234567890")
	require.NoError(t, err)
	require.Equal(t, true, status.Expired)
}
