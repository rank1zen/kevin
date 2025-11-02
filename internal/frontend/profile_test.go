package frontend_test

import (
	"context"
	"testing"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/stretchr/testify/require"
)

func TestProfileService_UpdateProfile(t *testing.T) {
	ctx := context.Background()

	service := (*frontend.ProfileService)(SetupHandler(ctx, t))

	err := service.UpdateProfile(ctx, frontend.UpdateProfileRequest{Name: "orrange", Tag: "NA1"})
	require.NoError(t, err)

	service.GetSummonerPage(ctx, frontend.GetSummonerPageRequest{})
}
