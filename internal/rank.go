package internal

import (
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

// RankStatus2 indicates the status of a summoner's rank.
type RankStatus2 struct {
	PUUID riot.PUUID

	// EffectiveDate indicates the time this status was taken.
	EffectiveDate time.Time

	// Detail is rank detail. A nil value indicates the summoner is
	// unranked.
	Detail *ZRankDetail
}

func NewRankStatus2(option RankStatus2Option) RankStatus2 {
	m := RankStatus2{}
	option(&m)
	return m
}

type RankStatus2Option func(*RankStatus2)

// Deprecated: get rid of.
type ZRankRecord struct {
	PUUID riot.PUUID

	// EffectiveDate indicates the time this record was taken.
	EffectiveDate time.Time

	// EndDate indicates the time this record is no longer current.
	EndDate *time.Time

	// IsCurrent indicates whether the record is current.
	IsCurrent bool

	// Detail is rank detail. A nil value indicates the summoner is
	// unranked.
	Detail *ZRankDetail
}

func ZNewRankRecord(from RankRecordFrom) (to ZRankRecord) {
	from(&to)
	return to
}

type RankRecordFrom func(*ZRankRecord)

// ZRankDetail contains details relating to a summoner's rank.
type ZRankDetail struct {
	Wins, Losses int

	Rank Rank
}

func ZNewRankDetail(from ZRankDetailOption) ZRankDetail {
	m := ZRankDetail{}
	from(&m)
	return m
}

type ZRankDetailOption func(*ZRankDetail)

func ZWithRiotLeagueEntry(rank riot.LeagueEntry) ZRankDetailOption {
	return func(m *ZRankDetail) {
		m.Wins = rank.Wins
		m.Losses = rank.Losses

		m.Rank = Rank{
			Tier:     rank.Tier,
			Division: rank.Division,
			LP:       rank.LeaguePoints,
		}
	}
}
