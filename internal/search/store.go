package search

import "context"

type Store interface {
	SearchSummoner(ctx context.Context, q string) ([]SearchResult, error)

	// SearchByNameTag searches for profiles by name and tag. It should return 10
	// profiles.
	SearchByNameTag(ctx context.Context, name, tag string) ([]Profile, error)
}
