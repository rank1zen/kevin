package ladder

import (
	"context"
	"fmt"

	"github.com/rank1zen/kevin/internal/riot"
)

type LeaderboardService struct {
	riot  *riot.Client
	store LeaderboardStore
}

type GetLeaderboardRequest struct {
	Region string `json:"region"`
}

// GetLeaderboard retrieves the leaderboard for ranked solo queue in a given region.
//
// NOTE: Only gets challenger league for now.
func (s *LeaderboardService) GetLeaderboard(ctx context.Context, req GetLeaderboardRequest) (*Leaderboard, error) {
	if req.Region == "" {
		req.Region = "NA1"
	}

	riotEntries, err := s.riot.League.GetChallengerLeague(ctx, req.Region, "RANKED_SOLO_5x5")
	if err != nil {
		return nil, fmt.Errorf("failed to get challenger league: %w", err)
	}

	result := Leaderboard{
		Region:  req.Region,
		Entries: []LeaderboardEntry{},
	}

	for _, entry := range riotEntries.Entries {
		result.Entries = append(result.Entries, LeaderboardEntry{
			PUUID:    entry.PUUID,
			Name:     "TODO: No name",
			Tag:      "TODO: No tag",
			Tier:     riotEntries.Tier,
			Division: entry.Rank,
			LP:       entry.LeaguePoints,
			Wins:     entry.Wins,
			Losses:   entry.Losses,
		})
	}

	return &result, nil
}
