package postgres

import (
	"context"

	"github.com/rank1zen/kevin/internal"
)

type LeaderboardStore Store

func (s *LeaderboardStore) GetLeaderboard(ctx context.Context, region string, filter internal.LeaderboardFilter) (*internal.Leaderboard, error) {
	result := internal.Leaderboard{
		Region:  region,
		Entries: []internal.LeaderboardEntry{},
	}

	return &result, nil
}
