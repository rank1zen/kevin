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

// TODO: I really don't think this MatchPage with days bucket is natural to the domain
type MatchPage struct {
	Days []MatchDay
}

type MatchDay struct {
	Day time.Time
}

// SummonerChampion is stats averaged by champion
type SummonerChampion struct {
	Champion               Champion
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
