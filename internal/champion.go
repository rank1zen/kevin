package internal

// ChampionStats is an aggregate average of a players stats on a specific
// champion over some number of games.
// TODO: rename to Champion
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

type ChampionStatsSeason struct {
	// TODO: stats for some season, for now just include everything

	List []ChampionStats
}
