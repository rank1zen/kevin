package internal

import (
	"context"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

// LiveMatchStore handles live match records.
type LiveMatchStore interface {
	CreateLiveMatch(ctx context.Context, match *LiveMatch) error

	GetLiveMatch(ctx context.Context, id string) (*LiveMatch, error)

	GetUserLiveMatch(ctx context.Context, puuid riot.PUUID) (*LiveMatch, error)
}

type LiveMatch struct {
	ID string

	// Date is game start timestamp
	Date time.Time

	// Participants are the players in this current match. There is no
	// chosen order.
	Participants [10]LiveParticipant
}

// LiveParticipant are currently in a match. NOTE: there are not a lot of
// fields are not available in an on-going game, including the summoners
// position.
type LiveParticipant struct {
	PUUID riot.PUUID

	MatchID string

	ChampionID int

	Runes RunePage

	SummonersIDs [2]int

	TeamID int
}
