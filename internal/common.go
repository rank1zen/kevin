package internal

type Season int

// NOTE: not fully implemented
const (
	Season2020 Season = 1
	SeasonAll  Season = -1
)

// PUUID is an encrypted riot PUUID. Exact length of 78 characters.
type PUUID string

func (id PUUID) String() string {
	return string(id)
}

// AccountID is an encrypted riot account ID. Max length of 56 characters.
type AccountID string

func (id AccountID) String() string {
	return string(id)
}

// SummonerID is an encrypted riot summoner ID. Max length of 63 characters.
type SummonerID string

func (id SummonerID) String() string {
	return string(id)
}

// ProfileIconID is the ID of the summoner profile icon.
type ProfileIconID int

type ItemID int

func (id ItemID) IconUrl() string {
	return ""
}

// ItemIDs are the 6 inventory slots and the 1 trinket slot.
// A value of nil means there is no item.
type ItemIDs [7]*ItemID

type SummsID int

func (id SummsID) IconUrl() string {
	return ""
}

// SummsIDs are the 2 summoner spells.
type SummsIDs [2]SummsID

// ChampionID is the ID of a champion.
type ChampionID int

func (id ChampionID) IconUrl() string {
	return ""
}

// MatchID is the ID of a match.
type MatchID string

func (id MatchID) String() string {
	return string(id)
}

// ParticipantID is the ID of a participant in a match.
type ParticipantID int

type TeamID int

// GameVersion is the patch a match was played on.
type GameVersion string

type RuneID int

func (id RuneID) IconUrl() string {
	return ""
}

type Runes struct {
	PrimaryTree     RuneID
	PrimaryKeystone RuneID
	PrimaryA        RuneID
	PrimaryB        RuneID
	PrimaryC        RuneID
	SecondaryTree   RuneID
	SecondaryA      RuneID
	SecondaryB      RuneID
	MiniOffense     RuneID
	MiniFlex        RuneID
	MiniDefense     RuneID
}

type RuneList [11]RuneID

func (r Runes) ToList() RuneList {
	return RuneList{
		r.PrimaryTree,
		r.PrimaryKeystone,
		r.PrimaryA,
		r.PrimaryB,
		r.PrimaryC,
		r.SecondaryTree,
		r.SecondaryA,
		r.SecondaryB,
		r.MiniOffense,
		r.MiniFlex,
		r.MiniDefense,
	}
}
