package internal

import (
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type RiotToProfileMapper struct {
	Account riot.Account

	Rank *riot.LeagueEntry

	EffectiveDate time.Time
}

func (m RiotToProfileMapper) Convert() Profile {
	profile := Profile{
		PUUID:   m.Account.PUUID,
		Name:    m.Account.GameName,
		Tagline: m.Account.TagLine,
		Rank: RankStatus{
			PUUID:         m.Account.PUUID,
			EffectiveDate: m.EffectiveDate,
			Detail:        nil,
		},
	}

	if m.Rank != nil {
		profile.Rank.Detail = &RankDetail{
			Wins:   m.Rank.Wins,
			Losses: m.Rank.Losses,
			Rank: Rank{
				Tier:     m.Rank.Tier,
				Division: m.Rank.Division,
				LP:       m.Rank.LeaguePoints,
			},
		}
	}

	return profile
}
