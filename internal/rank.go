package internal

import (
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

// RankDetail contains details relating to a summoner's rank.
type RankDetail struct {
	Wins, Losses int

	Rank Rank
}

// RankStatus indicates the status of a summoner's rank.
type RankStatus struct {
	PUUID riot.PUUID

	// EffectiveDate indicates the time this status was taken.
	EffectiveDate time.Time

	// Detail is rank detail. A nil value indicates the summoner is
	// unranked.
	Detail *RankDetail
}
