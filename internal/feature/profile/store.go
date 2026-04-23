package profile

import (
	"context"
	"time"
)

type ProfileStore interface {
	RecordProfile(ctx context.Context, profile *Profile) error

	GetProfile(ctx context.Context, puuid string) (*Profile, error)

	GetMatchlist(ctx context.Context, puuid string, start, end time.Time) ([]Match, error)

	GetChampions(ctx context.Context, puuid string, start, end time.Time) ([]ChampionAverage, error)

	GetRankHistory(ctx context.Context, puuid string, start, end time.Time) ([]RankStatus, error)
}
