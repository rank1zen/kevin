package internal_test

import (
	"context"
	"testing"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/require"
)

func TestDatasource_GetMatchDetail(t *testing.T) {
	ctx := context.Background()

	riotClient := riot.NewClient("")

	datasource := internal.NewDatasource(riotClient, nil)

	match, err := datasource.GetMatchDetail(ctx, riot.RegionNA1, "m")
	require.NoError(t, err)
}
