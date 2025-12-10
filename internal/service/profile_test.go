package service_test

import (
	"context"
	"testing"

	"github.com/rank1zen/kevin/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfileService_GetProfile(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	ds := SetupDatasource(ctx, t)

	req := service.GetProfileRequest{
		Name: "T1 OK GOOD YES",
		Tag:  "NA1",
	}

	profile, err := (*service.ProfileService)(ds).GetProfile(ctx, req)
	require.NoError(t, err)

	assert.EqualValues(t, profile.PUUID, "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")
}
