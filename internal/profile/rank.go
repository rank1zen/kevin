package profile

import "time"

// RankDetail contains detailed rank information for a summoner.
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

type Rank struct {
	Tier         string
	Rank         string
	LeaguePoints int
}
