package riotmapper

import (
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

type Profile struct {
	Account       riot.Account
	Rank          *riot.LeagueEntry
	EffectiveDate time.Time
}

func MapProfile(profile *Profile) *internal.Profile {
	result := internal.Profile{
		PUUID:   profile.Account.PUUID,
		Name:    profile.Account.GameName,
		Tagline: profile.Account.TagLine,
		Rank: internal.RankStatus{
			PUUID:         profile.Account.PUUID,
			EffectiveDate: profile.EffectiveDate,
			Detail:        nil,
		},
	}

	if profile.Rank != nil {
		result.Rank.Detail = &internal.RankDetail{
			Wins:   profile.Rank.Wins,
			Losses: profile.Rank.Losses,
			Rank: internal.Rank{
				Tier:     profile.Rank.Tier,
				Division: profile.Rank.Division,
				LP:       profile.Rank.LeaguePoints,
			},
		}
	}

	return &result
}
