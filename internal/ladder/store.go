package ladder

import "context"

// LeaderboardStore provides access to leaderboard data.
type LeaderboardStore interface {
	GetLeaderboard(ctx context.Context, region string, filter LeaderboardFilter) (*Leaderboard, error)
}
