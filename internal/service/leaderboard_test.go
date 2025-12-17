package service_test

import (
	"context"
	"testing"

	"github.com/rank1zen/kevin/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLeaderboardService_GetLeaderboard
func TestLeaderboardService_GetLeaderboard(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	ds := SetupDatasource(ctx, t)

	req := service.GetLeaderboardRequest{}

	result, err := (*service.LeaderboardService)(ds).GetLeaderboard(ctx, req)
	require.NoError(t, err)

	assert.Equal(t, "NA1", result.Region)
	assert.Greater(t, len(result.Entries), 0)
}
