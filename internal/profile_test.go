package internal_test

import (
	"context"
	"testing"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfileService_GetProfile(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	ds := SetupDatasource(ctx, t)

	req := internal.GetProfileRequest{
		Name: "T1 OK GOOD YES",
		Tag:  "NA1",
	}

	profile, err := (*internal.ProfileService)(ds).GetProfile(ctx, req)
	require.NoError(t, err)

	assert.Equal(t, riot.RegionNA1, *req.Region)
	assert.EqualValues(t, profile.PUUID, "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")
}
