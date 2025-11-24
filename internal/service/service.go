package service

import (
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
	match   internal.MatchStore
	profile internal.ProfileStore

	riot *riot.Client
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
