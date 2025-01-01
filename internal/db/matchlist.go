package db

import (
	"context"

	"github.com/rank1zen/yujin/internal"
)

func (db *DB) GetMatchList(ctx context.Context, puuid internal.PUUID, page int, ensure bool) ([]internal.RiotMatchParticipant, error) {
	panic("deprecated: remove this")
}

func (db *DB) CheckMatchIDs(context.Context, []internal.MatchID) ([]internal.MatchID, error) {
	panic("deprecated: remove this")
}
