package match

import (
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

// Match represents a record of a ranked match.
type Match struct {
	ID           string
	Date         time.Time
	Duration     time.Duration
	Version      string
	WinnerID     int
	Participants [10]Participant
}

// Participant represents a record of a summoner in a ranked match.
type Participant struct {
	PUUID                  riot.PUUID
	MatchID                string
	TeamID                 int
	ChampionID             int
	ChampionLevel          int
	TeamPosition           internal.TeamPosition
	SummonerIDs            [2]int
	Runes                  internal.RunePage
	Items                  [7]int
	Kills, Deaths, Assists int
	KillParticipation      float32
	CreepScore             int
	CreepScorePerMinute    float32
	DamageDealt            int
	DamageTaken            int
	DamageDeltaEnemy       int
	DamagePercentageTeam   float32
	GoldEarned             int
	GoldDeltaEnemy         int
	GoldPercentageTeam     float32
	VisionScore            int
	PinkWardsBought        int
}

// MatchDetail are details relating to a match record.
type MatchDetail struct {
	// ID is a region+number, which forms an identifier. NOTE: should
	// switch to new match ID type.
	ID string

	// Date is the end timestamp of the match.
	Date time.Time

	// Duration is the length of the match.
	Duration time.Duration

	// Version is the game version.
	Version string

	// WinnerID is the ID of the winning team. NOTE: should switch to new
	// TeamID type.
	WinnerID int

	// Participants are the players in this match. There is no chosen
	// order.
	Participants [10]ParticipantDetail
}

// ParticipantDetail is the details relating to a participant record.
type ParticipantDetail struct {
	PUUID                  riot.PUUID
	MatchID                string
	TeamID                 int
	ChampionID             int
	ChampionLevel          int
	TeamPosition           internal.TeamPosition
	SummonerIDs            [2]int
	Runes                  internal.RunePage
	Items                  [7]int
	Kills, Deaths, Assists int
	KillParticipation      float32
	CreepScore             int
	CreepScorePerMinute    float32
	DamageDealt            int
	DamageTaken            int
	DamageDeltaEnemy       int
	DamagePercentageTeam   float32
	GoldEarned             int
	GoldDeltaEnemy         int
	GoldPercentageTeam     float32
	VisionScore            int
	PinkWardsBought        int

	Name, Tag string

	// CurrentRank is the player's rank at the time of the match.
	CurrentRank *RankDetail
}

type RankDetail struct {
	Tier     string // e.g. Diamond
	Division string // e.g. III
	LP       int
}
