package frontend_test

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfileHandler_GetMatchHistory(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	ds := SetupDatasource(ctx, t)
	handler := frontend.ProfileHandler{Datasource: ds}

	tz, err := time.LoadLocation("America/Toronto")
	require.NoError(t, err)

	req := frontend.MatchHistoryRequest{
		Region:  riot.RegionNA1,
		PUUID:   internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg"),
		StartTS: time.Date(2025, time.August, 26, 22, 0, 0, 0, tz),
		EndTS:   time.Date(2025, time.August, 27, 0, 0, 0, 0, tz),
	}

	v, err := handler.GetMatchHistory(ctx, req)
	require.NoError(t, err)

	ids := []string{}
	for _, match := range v.Matchlist {
		ids = append(ids, match.MatchID)
	}

	assert.Equal(t, []string{"NA1_5356002067", "NA1_5355962057"}, ids)
}
