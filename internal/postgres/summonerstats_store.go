package postgres

import (
	"context"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

type SummonerStatsStore Store

func (db *SummonerStatsStore) GetChampions(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]internal.SummonerChampion, error) {
	matchStore := MatchStore{Tx: db.Pool}

	return matchStore.GetSummonerChampions(ctx, puuid, start, end)
}
