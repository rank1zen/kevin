package internal

import (
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type Rank struct {
	Tier   riot.Tier
	Division   riot.Division
	LP     int
}

type RankStatus struct {
	PUUID         string
	EffectiveDate time.Time
	Detail        *RankDetail
}

type RankRecord struct {
	PUUID         string
	EffectiveDate time.Time
	EndDate       *time.Time
	IsCurrent     bool
	Detail        *RankDetail
}

type RankDetail struct {
	Wins   int
	Losses int
	Tier   riot.Tier
	Division   riot.Division
	LP     int
}

func NewRankDetail(opts ...RankDetailOption) RankDetail {
	var m RankDetail
	for _, f := range opts {
		if err := f(&m); err != nil {
			panic(err)
		}
	}
	return m
}

type RankDetailOption func(*RankDetail) error

func WithRiotLeagueEntry(rank riot.LeagueEntry) RankDetailOption {
	return func(m *RankDetail) error {
		m.Wins = rank.Wins
		m.Losses = rank.Losses
		m.Tier = rank.Tier
		m.Division = rank.Division
		m.LP = rank.LeaguePoints
		return nil
	}
}
