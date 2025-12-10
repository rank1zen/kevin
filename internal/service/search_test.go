package service_test

import (
	"context"
	"testing"

	"github.com/rank1zen/kevin/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchService_SearchProfile(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	ds := SetupDatasource(ctx, t)

	err := (*service.ProfileService)(ds).UpdateProfile(ctx, service.UpdateProfileRequest{
		Name: "orrange",
		Tag:  "NA1",
	})
	require.NoError(t, err)

	req := service.SearchProfileRequest{
		Query: "orrange#NA1",
	}

	result, err := (*service.SearchService)(ds).SearchProfile(ctx, req)
	if assert.NoError(t, err) {
		assert.Equal(t, "orrange", result.Name)
		assert.Equal(t, "NA1", result.Tag)
		assert.Greater(t, len(result.Profiles), 0)
	}
}
