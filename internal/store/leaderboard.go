package postgres

import (
	"context"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
)

type LeaderboardStore Store

func (s *LeaderboardStore) GetLeaderboard(ctx context.Context, region string, filter internal.LeaderboardFilter) (*internal.Leaderboard, error) {
	var (
		rankStore     = postgres.RankStore{Tx: s.Pool}
		summonerStore = postgres.SummonerStore{Tx: s.Pool}
	)

	ranklist, err := rankStore.ListRanks(ctx, region, postgres.LeaderBoardOption{
		Start: filter.Start,
		Count: filter.Count,
	})

	if err != nil {
		return nil, err
	}

	result := internal.Leaderboard{
		Region:  region,
		Entries: []internal.LeaderboardEntry{},
	}

	for _, rank := range ranklist {
		summoner, err := summonerStore.GetSummoner(ctx, riot.PUUID(rank.Status.PUUID))
		if err != nil {
			return nil, err
		}

		result.Entries = append(result.Entries, internal.LeaderboardEntry{
			Region:   region,
			PUUID:    rank.Status.PUUID,
			Name:     summoner.Name,
			Tag:      summoner.Tagline,
			Tier:     rank.Detail.Tier,
			Division: rank.Detail.Division,
			LP:       rank.Detail.LP,
			Wins:     rank.Detail.Wins,
			Losses:   rank.Detail.Losses,
		})
	}

	return &result, nil
}
