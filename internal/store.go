package internal

import "context"

// Store represents a collection of data stores.
type Store struct {
	Profile       ProfileStore
	Match         MatchStore
	SummonerStats SummonerStatsStore
}

type DB interface {
	// Ping returns an error if the database is not reachable.
	Ping(ctx context.Context) error
}
