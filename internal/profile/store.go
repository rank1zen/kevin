package profile

import (
	"context"
	"time"
)

type Store interface {
	RecordProfile(ctx context.Context, profile *Profile) error

	GetProfile(ctx context.Context, puuid string) (*Profile, error)

	GetMatchlist(ctx context.Context, puuid string, start, end time.Time) ([]Match, error)

	GetChampions(ctx context.Context, puuid string, start, end time.Time) ([]ChampionAggregate, error)

	GetRankHistory(ctx context.Context, puuid string, start, end time.Time) ([]RankStatus, error)
}

type Stores struct{}

func (s Stores) RecordProfile(ctx context.Context, profile *Profile) error {
	//TODO implement me
	panic("implement me")
}

func (s Stores) GetProfile(ctx context.Context, puuid string) (*Profile, error) {
	//TODO implement me
	panic("implement me")
}

func (s Stores) GetMatchlist(ctx context.Context, puuid string, start, end time.Time) ([]Match, error) {
	//TODO implement me
	panic("implement me")
}

func (s Stores) GetChampions(ctx context.Context, puuid string, start, end time.Time) ([]ChampionAggregate, error) {
	//TODO implement me
	panic("implement me")
}

func (s Stores) GetRankHistory(ctx context.Context, puuid string, start, end time.Time) ([]RankStatus, error) {
	//TODO implement me
	panic("implement me")
}
