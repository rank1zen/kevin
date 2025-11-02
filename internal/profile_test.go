package internal_test

import (
	"context"
	"testing"

	"github.com/rank1zen/kevin/internal"
	"github.com/stretchr/testify/assert"
)

func TestProfileService_GetProfile(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	ds := SetupDatasource(ctx, t)
	service := (*internal.ProfileService)(ds)

	req := internal.GetProfileRequest{
		Name: "T1 OK GOOD YES",
		Tag:  "NA1",
	}

	_, err := service.GetProfile(ctx, req)
	assert.NoError(t, err)
}
