package internal

import (
	"time"
)

type Summoner struct {
	PUUID         string
	Name, Tagline string
	Platform      string
	SummonerID    string
}

type SummonerMatch struct {
	Date     time.Time
	Duration time.Duration
	LpDelta  *int
	Win      bool

	Participant
}

// SummonerChampion is a summoner's champion stats average over GamesPlayed.
type SummonerChampion struct {
	PUUID                  string
	GamesPlayed            int
	Wins, Losses           int
	Champion               Champion
	Kills, Deaths, Assists float32
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
