package internal

import (
	"github.com/rank1zen/kevin/internal/riot"
)

// Datasource manages interaction between the riot API and an internal store.
//
// Region parameters specify the region to search.
//
// TODO: region parameters not
// implemented.
//
// TODO: Datasource should be able to decide when to call the riot API, and
// when to use cache. probably want to cache something.
type Datasource struct {
	match   MatchStore
	profile ProfileStore

	riot *riot.Client
}

type Store interface {
	MatchStore() MatchStore
	ProfileStore() ProfileStore
}

func NewDatasource(client *riot.Client, store Store) *Datasource {
	return &Datasource{
		riot:    client,
		match:   store.MatchStore(),
		profile: store.ProfileStore(),
	}
}
