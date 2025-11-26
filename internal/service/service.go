package service

import (
	"context"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

// Service manages interaction between the riot API and an internal store.
//
// Region parameters specify the region to search.
//
// TODO: region parameters not
// implemented.
//
// TODO: Service should be able to decide when to call the riot API, and
// when to use cache. probably want to cache something.
type Service struct {
	store         internal.Store
	match         internal.MatchStore
	profile       internal.ProfileStore
	summonerstats internal.SummonerStatsStore

	riot *riot.Client
}

func (s *Service) CheckHealth(ctx context.Context) error {
	if _, err := s.riot.Account.GetAccountByRiotID(ctx, riot.RegionNA1, "orrange", "NA1"); err != nil {
		return err
	}

	if err := s.store.Ping(ctx); err != nil {
		return err
	}

	return nil
}

type Store interface {
	MatchStore() internal.MatchStore
	ProfileStore() internal.ProfileStore
}

func NewService(client *riot.Client, store Store) *Service {
	return &Service{
		riot:    client,
		match:   store.MatchStore(),
		profile: store.ProfileStore(),
	}
}
