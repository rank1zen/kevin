package internal

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rank1zen/kevin/internal/riot"
)

type RankDetail struct {
	Wins   int
	Losses int
	Tier   riot.Tier
	Rank   riot.Rank
	LP     int
}

type RankStatus struct {
	PUUID         string
	EffectiveDate time.Time
	EndDate       pgtype.Date
	IsCurrent     bool
	IsRanked      bool
	Detail        *RankDetail
}

func NewRankStatus(opts ...RankStatusOption) RankStatus {
	var rs RankStatus
	for _, f := range opts {
		if err := f(&rs); err != nil {
			panic(err)
		}
	}
	return rs
}

type RankStatusOption func(*RankStatus) error

func WithRiotLeagueEntry(puuid string, t time.Time, rank *riot.LeagueEntry) RankStatusOption {
	return func(rs *RankStatus) error {
		rs.PUUID = puuid
		rs.EffectiveDate = t
		rs.EndDate = pgtype.Date{InfinityModifier: 1, Valid: true}
		rs.IsCurrent = true
		if rank == nil {
			rs.Detail = nil
			rs.IsRanked = false
		} else {
			rs.Detail = &RankDetail{
				Wins:   rank.Wins,
				Losses: rank.Losses,
				Tier:   rank.Tier,
				Rank:   rank.Rank,
				LP:     rank.LeaguePoints,
			}
		}
		return nil
	}
}
