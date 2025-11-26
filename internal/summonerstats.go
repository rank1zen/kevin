package internal

import (
	"context"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

// SummonerStatsStore defines operations for retrieving aggregated summoner statistics.
type SummonerStatsStore interface {
	GetChampions(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]SummonerChampion, error)
}

// SummonerChampion is a summoner's champion stats averaged over GamesPlayed.
type SummonerChampion struct {
	PUUID riot.PUUID

	// NOTE: Champion type should be specified by ddragon package.
	Champion Champion

	GamesPlayed int

	Wins, Losses int

	AverageKillsPerGame float32

	AverageDeathsPerGame float32

	AverageAssistsPerGame float32

	AverageKillParticipationPerGame float32

	AverageCreepScorePerGame float32

	AverageCreepScorePerMinutePerGame float32

	AverageDamageDealtPerGame float32

	AverageDamageTakenPerGame float32

	AverageDamageDeltaEnemyPerGame float32

	AverageDamagePercentagePerGame float32

	AverageGoldEarnedPerGame float32

	AverageGoldDeltaEnemyPerGame float32

	AverageGoldPercentagePerGame float32

	AverageVisionScorePerGame float32

	AveragePinkWardsBoughtPerGame float32
}
