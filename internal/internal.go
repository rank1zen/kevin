package internal

import (
	"context"
	"fmt"
	"time"
)

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

type LiveMatchParticipant struct {
	Puuid          PUUID
	Date           time.Time
	TeamID         TeamID
	SummonerID     SummonerID
	Champion       ChampionID
	Summoners      SummsIDs
	Runes          Runes
	BannedChampion *ChampionID
}

type LiveMatchParticipantList [10]LiveMatchParticipant

type LiveMatchTeamList [5]LiveMatchParticipant

type LiveMatch struct {
	StartTimestamp time.Time
	Length         time.Duration
	IDs            [10]PUUID
	Participant    LiveMatchParticipantList
}

// TODO: we should further implement checking the order of the list
func (m *LiveMatch) GetParticipants() LiveMatchParticipantList {
	return m.Participant
}

// Match is strictly a ranked, soloq, 5v5, match on summoners rift.
type Match struct {
	ID              MatchID
	DataVersion     string
	Patch           GameVersion
	CreateTimestamp time.Time
	StartTimestamp  time.Time
	EndTimestamp    time.Time
	Duration        time.Duration
	EndOfGameResult string
	Participants    MatchParticipantList
}

type MatchParticipantList [10]MatchParticipant

type MatchTeamList [5]MatchParticipant

func (m *Match) GetParticipants() MatchParticipantList {
	return m.Participants
}

func (m *Match) GetTeams() [2]MatchTeam {
	return [2]MatchTeam{}
}

type MatchParticipant struct {
	ID       ParticipantID
	Puuid    PUUID
	Match    MatchID
	Team     TeamID
	Summoner SummonerID

	SummonerLevel             int
	SummonerName              string
	RiotIDGameName            string
	RiotIDName                string
	RiotIDTagline             string
	ChampionLevel             int
	ChampionID                ChampionID
	ChampionName              string
	GameEndedInEarlySurrender bool
	GameEndedInSurrender      bool
	Items                     ItemIDs
	Runes                     Runes
	Role                      string
	Summoners                 SummsIDs
	TeamEarlySurrendered      bool
	TeamPosition              string
	TimePlayed                int
	Win                       bool

	// Post game stats

	Assists                        int
	DamageDealtToBuildings         int
	DamageDealtToObjectives        int
	DamageDealtToTurrets           int
	DamageSelfMitigated            int
	Deaths                         int
	DetectorWardsPlaced            int
	FirstBloodAssist               bool
	FirstBloodKill                 bool
	FirstTowerAssist               bool
	FirstTowerKill                 bool
	GoldEarned                     int
	GoldSpent                      int
	IndividualPosition             string
	InhibitorKills                 int
	InhibitorTakedowns             int
	InhibitorsLost                 int
	Kills                          int
	MagicDamageDealt               int
	MagicDamageDealtToChampions    int
	MagicDamageTaken               int
	PhysicalDamageDealt            int
	PhysicalDamageDealtToChampions int
	PhysicalDamageTaken            int
	SightWardsBoughtInGame         int
	TotalDamageDealt               int
	TotalDamageDealtToChampions    int
	TotalDamageShieldedOnTeammates int
	TotalDamageTaken               int
	TotalHeal                      int
	TotalHealsOnTeammates          int
	TotalMinionsKilled             int
	TrueDamageDealt                int
	TrueDamageDealtToChampions     int
	TrueDamageTaken                int
	VisionScore                    int
	VisionWardsBoughtInGame        int
	WardsKilled                    int
	WardsPlaced                    int
	NeutralMinionsKilled           int
}

type MatchTeam struct {
	ID      TeamID
	MatchID MatchID
	Win     bool
}

// Profile represents a record of a summoner.
type Profile struct {
	RecordDate    time.Time
	RevisionDate  time.Time
	Puuid         PUUID
	SummonerID    SummonerID
	AccountID     AccountID
	Name          string
	Tagline       string
	Level         int
	ProfileIconID ProfileIconID
	Rank          *RankRecord
}

type SummonerRecord struct {
	ID            SummonerID
	Puuid         PUUID
	RevisionDate  time.Time
	ValidFrom     time.Time
	ValidTo       time.Time
	EnteredAt     time.Time
	Level         int
	ProfileIconID ProfileIconID
}

type RankRecord struct {
	Puuid     PUUID
	Timestamp time.Time
	LeagueID  string
	Wins      int
	Losses    int
	Tier      string
	Division  string
	LP        int
}

// TODO: implement the cases for Challenger, GM, Masters, and Unranked.
func (r RankRecord) String() string {
	return fmt.Sprintf("%s %s %d", r.Tier, r.Division, r.LP)
}

type RankSnapshot struct {
	RankRecord

	LpDelta *int
}

// ChampionStats is an aggregate average of a players stats on a specific
// champion over some time interval.
type ChampionStats struct {
	Puuid             PUUID
	Champion          ChampionID
	GamesPlayed       int
	WinPercentage     float32
	Wins              int
	Losses            int
	LpDelta           int
	Kills             float32
	Deaths            float32
	Assists           float32
	KillParticipation float32
	CreepScore        float32
	CsPerMinute       float32
	Damage            float32
	DamagePercentage  float32
	DamageDelta       float32
	GoldEarned        float32
	GoldPercentage    float32
	GoldDelta         float32
	VisionScore       float32
}

type Repository interface {
	CheckProfileExists(context.Context, PUUID) (bool, error)

	UpdateProfile(context.Context, Profile) error

	GetProfile(context.Context, PUUID) (Profile, error)

	GetRankList(context.Context, PUUID) ([]RankRecord, error)

	GetMatchList(context.Context, PUUID, int, bool) ([]MatchParticipant, error)

	GetChampionList(context.Context, PUUID) ([]ChampionStats, error)

	CreateMatch(context.Context, Match) error
}

type RiotClient interface {
	GetProfile(context.Context, PUUID) (Profile, error)

	GetLiveMatch(context.Context, PUUID) (LiveMatch, error)

	GetMatchList(context.Context, PUUID) ([]MatchID, error)

	GetMatch(context.Context, MatchID) (Match, error)
}
