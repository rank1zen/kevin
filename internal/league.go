package internal

import "github.com/rank1zen/kevin/internal/riot"

// Rank is summoner rank.
type Rank struct {
	// e.g. Diamond
	Tier riot.Tier

	// e.g. III
	Division riot.Division

	LP int
}

type Champion int

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

type RunePageOption func(*RunePage) error

func NewRunePage(opts ...RunePageOption) (runes RunePage) {
	for _, f := range opts {
		f(&runes)
	}
	return runes
}

func WithRiotPerks(perks *riot.MatchPerks) RunePageOption {
	return func(runes *RunePage) error {
		runes.PrimaryTree = perks.Styles[0].Style
		runes.PrimaryKeystone = perks.Styles[0].Selections[0].Perk
		runes.PrimaryA = perks.Styles[0].Selections[1].Perk
		runes.PrimaryB = perks.Styles[0].Selections[2].Perk
		runes.PrimaryC = perks.Styles[0].Selections[3].Perk
		runes.SecondaryTree = perks.Styles[1].Style
		runes.SecondaryA = perks.Styles[1].Selections[0].Perk
		runes.SecondaryB = perks.Styles[1].Selections[1].Perk
		runes.MiniOffense = perks.StatPerks.Offense
		runes.MiniFlex = perks.StatPerks.Flex
		runes.MiniDefense = perks.StatPerks.Defense
		return nil
	}
}

func WithIntList(ids [11]int) RunePageOption {
	return func(runes *RunePage) error {
		runes.PrimaryTree = ids[0]
		runes.PrimaryKeystone = ids[1]
		runes.PrimaryA = ids[2]
		runes.PrimaryB = ids[3]
		runes.PrimaryC = ids[4]
		runes.SecondaryTree = ids[5]
		runes.SecondaryA = ids[6]
		runes.SecondaryB = ids[7]
		runes.MiniOffense = ids[8]
		runes.MiniFlex = ids[9]
		runes.MiniDefense = ids[10]
		return nil
	}
}

func WithRiotSpectatorPerks(perks *riot.LivePerks) RunePageOption {
	return func(runes *RunePage) error {
		runes.PrimaryTree = perks.PerkStyle
		runes.PrimaryKeystone = perks.PerkIDs[0]
		runes.PrimaryA = perks.PerkIDs[1]
		runes.PrimaryB = perks.PerkIDs[2]
		runes.PrimaryC = perks.PerkIDs[3]
		runes.SecondaryTree = perks.PerkSubStyle
		runes.SecondaryA = perks.PerkIDs[4]
		runes.SecondaryB = perks.PerkIDs[5]
		runes.MiniOffense = perks.PerkIDs[6]
		runes.MiniFlex = perks.PerkIDs[7]
		runes.MiniDefense = perks.PerkIDs[8]
		return nil
	}
}
