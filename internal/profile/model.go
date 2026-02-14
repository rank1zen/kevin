package profile

import (
	"time"

	"github.com/rank1zen/kevin/internal"
)

// Profile is a summoner's profile. PUUID is unique and immutable. Name +
// Tagline is unique but mutable.
type Profile struct {
	PUUID         string
	Name, Tagline string

	// Rank is the summoners most recent rank record in store, meaning it could be
	// out-of-date.
	Rank RankStatus
}

type RankDetail struct {
	Wins, Losses int
	Tier         string // e.g. Diamond
	Division     string // e.g. III
	LP           int
}

// RankStatus indicates the status of a summoner's rank.
type RankStatus struct {
	PUUID string

	// EffectiveDate indicates the time this status was taken.
	EffectiveDate time.Time

	// Detail is rank detail. A nil value indicates the summoner is
	// unranked.
	Detail *RankDetail
}

type Match struct {
	PUUID                  string
	MatchID                string
	TeamID                 int
	ChampionID             int
	ChampionLevel          int
	TeamPosition           string
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

	Date     time.Time
	Duration time.Duration
	Win      bool

	// RankBefore is the summoner's rank just before the match. A nil value
	// indicates this no record was taken.
	RankBefore *RankStatus
	// RankBefore is the summoner's rank just after the match. A nil value
	// indicates this no record was taken.
	RankAfter *RankStatus
}

// ChampionAverage is a summoner's champion stats averaged over GamesPlayed.
type ChampionAverage struct {
	PUUID                             string
	Champion                          int
	GamesPlayed                       int
	Wins, Losses                      int
	AverageKillsPerGame               float32
	AverageDeathsPerGame              float32
	AverageAssistsPerGame             float32
	AverageKillParticipationPerGame   float32
	AverageCreepScorePerGame          float32
	AverageCreepScorePerMinutePerGame float32
	AverageDamageDealtPerGame         float32
	AverageDamageTakenPerGame         float32
	AverageDamageDeltaEnemyPerGame    float32
	AverageDamagePercentagePerGame    float32
	AverageGoldEarnedPerGame          float32
	AverageGoldDeltaEnemyPerGame      float32
	AverageGoldPercentagePerGame      float32
	AverageVisionScorePerGame         float32
	AveragePinkWardsBoughtPerGame     float32
}
