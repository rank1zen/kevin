package profile

type ChampionAggregate struct {
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
