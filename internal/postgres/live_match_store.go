package postgres

import (
	"context"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

type LiveMatchStore Store

func (db *LiveMatchStore) CreateLiveMatch(ctx context.Context, match *internal.LiveMatch) error {
	panic("not implemented")
}

func (db *LiveMatchStore) GetLiveMatch(ctx context.Context, id string) (*internal.LiveMatch, error) {
	panic("not implemented")
}

func (db *LiveMatchStore) GetUserLiveMatch(ctx context.Context, puuid riot.PUUID) (*internal.LiveMatch, error) {
	panic("not implemented")
}
