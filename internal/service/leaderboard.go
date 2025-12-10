package service

import (
	"context"
	"fmt"
)

type LeaderboardService Service

// GetLeaderboard retrieves the leaderboard for ranked solo queue in a given region.
//
// NOTE: Only gets challenger league for now.
func (s *LeaderboardService) GetLeaderboard(ctx context.Context, region string) ([]LeaderboardEntry, error) {
	riotEntries, err := s.riot.League.GetChallengerLeague(ctx, "NA1", "RANKED_SOLO_5x5")
	if err != nil {
		return nil, fmt.Errorf("failed to get challenger league: %w", err)
	}

	result := []LeaderboardEntry{}

	for _, entry := range riotEntries.Entries {
		result = append(result, LeaderboardEntry{
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

	return result, nil
}

// LeaderboardEntry is an entry for a player in a leaderboard.
type LeaderboardEntry struct {
	PUUID string

	Name, Tag string

	Tier string

	Division string

	LP int

	Wins, Losses int
}
