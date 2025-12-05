package service_test

import (
	"context"
	"testing"

	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLiveMatchService_GetLiveMatchByID(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	ds := SetupDatasource(ctx, t)

	req := service.GetLiveMatchByIDRequest{
		MatchID: "NA1_5346312088",
	}

	match, err := (*service.LiveMatchService)(ds).GetLiveMatchByID(ctx, req)
	require.NoError(t, err)

	assert.Equal(t, "NA1_5346312088", match.ID)
	assert.Equal(t, riot.RegionNA1, match.Region)
}
