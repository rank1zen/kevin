package postgres

import (
	"context"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type LiveMatchObjects struct{ Tx Tx }

func (db *LiveMatchObjects) Get(ctx context.Context, id string) (*LiveMatchStatus, error) {
	panic("not implemented")
}

func (db *LiveMatchObjects) GetByPUUID(ctx context.Context, puuid string) (*LiveMatchStatus, error) {
	panic("not implemented")
}

type LiveMatch struct {
	ID   string    `db:"match_id"`
	Date time.Time `db:"date"`
}

type LiveParticipant struct {
	PUUID        riot.PUUID `db:"puuid"`
	MatchID      string     `db:"match_id"`
	ChampionID   int        `db:"champion_id"`
	Runes        [11]int    `db:"runes"`
	SummonersIDs [2]int     `db:"summoners"`
	TeamID       int        `db:"team"`
}
