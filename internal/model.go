package internal

import (
	"github.com/rank1zen/kevin/internal/riot"
)

// Rank is summoner rank.
type Rank struct {
	// e.g. Diamond
	Tier riot.Tier

	// e.g. III
	Division riot.Division

	LP int
}

// TeamPosition is a players position in a game. There must be one of each on
// a team: TOP, JNG, MID, BOT, SUP.
type TeamPosition int

const (
	TeamPositionTop TeamPosition = iota
	TeamPositionJungle
	TeamPositionMiddle
	TeamPositionBottom
	TeamPositionSupport
)

// RunePage is the set of runes chosen by a summoner in a match.
type RunePage struct {
	PrimaryTree     int
	PrimaryKeystone int
	PrimaryA        int
	PrimaryB        int
	PrimaryC        int
	SecondaryTree   int
	SecondaryA      int
	SecondaryB      int
	MiniOffense     int
	MiniFlex        int
	MiniDefense     int
}
