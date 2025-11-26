package internal

import (
	"context"

	"github.com/rank1zen/kevin/internal/riot"
)

type ProfileStore interface {
	RecordProfile(ctx context.Context, profile *Profile) error

	GetProfile(ctx context.Context, puuid riot.PUUID) (*Profile, error)

	SearchSummoner(ctx context.Context, q string) ([]SearchResult, error)
}

type Profile struct {
	PUUID riot.PUUID

	Name, Tagline string

	Rank RankStatus
}

type SearchResult struct {
	PUUID riot.PUUID

	Name, Tagline string

	// Rank is the summoners most recent rank record in store.
	Rank *RankStatus
}
