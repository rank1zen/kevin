package profile

import "time"

// Match represents a completed match for a summoner.
type Match struct {
	PUUID                  string
	MatchID                string
	TeamID                 int
	ChampionID             int
	ChampionLevel          int
	TeamPosition           string
	SummonerIDs            [2]int
	Runes                  RuneTree
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
