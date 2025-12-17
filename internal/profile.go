package internal

import (
	"context"

	"github.com/rank1zen/kevin/internal/riot"
)

type ProfileStore interface {
	RecordProfile(ctx context.Context, profile *Profile) error

	GetProfile(ctx context.Context, puuid riot.PUUID) (*Profile, error)

	SearchSummoner(ctx context.Context, q string) ([]SearchResult, error)

	// SearchByNameTag searches for profiles by name and tag. It should return 10
	// profiles.
	SearchByNameTag(ctx context.Context, name, tag string) ([]Profile, error)
}

// Profile is a summoner's profile. PUUID is unique and immutable. Name +
// Tagline is unique but mutable.
type Profile struct {
	PUUID riot.PUUID

	Name, Tagline string

	// Rank is the summoners most recent rank record in store, meaning it could be
	// out-of-date.
	Rank RankStatus
}

type SearchResult struct {
	PUUID riot.PUUID

	Name, Tagline string

	// Rank is the summoners most recent rank record in store.
	Rank *RankStatus
}
