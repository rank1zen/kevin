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

func TestLeaderboardStore_GetLeaderboard(t *testing.T) {
	ctx := context.Background()
	db := setupStore(ctx, t)

	var (
		leaderboardStore = (*store.LeaderboardStore)(db)
		profileStore     = (*store.ProfileStore)(db)
	)

	err := profileStore.RecordProfile(ctx, &internal.Profile{
		PUUID:   "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg",
		Name:    "T1 OK GOOD YES",
		Tagline: "NA1",
		Rank: internal.RankStatus{
			PUUID:         "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg",
			EffectiveDate: time.Now(),
			Detail: &internal.RankDetail{
				Wins:   12,
				Losses: 12,
				Rank: internal.Rank{
					Tier:     riot.TierDiamond,
					Division: riot.Division3,
					LP:       12,
				},
			},
		},
	})
	require.NoError(t, err)

	err = profileStore.RecordProfile(ctx, &internal.Profile{
		PUUID:   "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
		Name:    "orrange",
		Tagline: "NA1",
		Rank: internal.RankStatus{
			PUUID:         "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
			EffectiveDate: time.Now(),
			Detail: &internal.RankDetail{
				Wins:   14,
				Losses: 14,
				Rank: internal.Rank{
					Tier:     riot.TierDiamond,
					Division: riot.Division3,
					LP:       13,
				},
			},
		},
	})
	require.NoError(t, err)

	leaderboard, err := leaderboardStore.GetLeaderboard(ctx, "NA1", internal.LeaderboardFilter{
		Start: 0,
		Count: 10,
	})
	require.NoError(t, err)
	require.Len(t, leaderboard.Entries, 2)

	assert.Equal(t, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q", leaderboard.Entries[0].PUUID)
}
