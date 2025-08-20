package postgres

import (
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

type PostgresToSummonerMatchMapper struct {
	Match Match

	Participant Participant

	BeforeStatus *RankStatus
	BeforeDetail *RankDetail

	AfterStatus *RankStatus
	AfterDetail *RankDetail
}

func (mapper PostgresToSummonerMatchMapper) Map() internal.SummonerMatch {
	p := mapper.Participant
	m := mapper.Match

	win := false
	if p.TeamID == m.WinnerID {
		win = true
	}

	result := internal.SummonerMatch{
		Participant: internal.Participant{
			PUUID:                internal.NewPUUIDFromString(p.PUUID),
			MatchID:              p.MatchID,
			TeamID:               p.TeamID,
			ChampionID:           p.ChampionID,
			ChampionLevel:        p.ChampionLevel,
			TeamPosition:         convertStringToTeamPosition(p.TeamPosition),
			SummonerIDs:          p.SummonerIDs,
			Runes:                convertListToRunePage(p.Runes),
			Items:                p.Items,
			Kills:                p.Kills,
			Deaths:               p.Deaths,
			Assists:              p.Assists,
			KillParticipation:    p.KillParticipation,
			CreepScore:           p.CreepScore,
			CreepScorePerMinute:  p.CreepScorePerMinute,
			DamageDealt:          p.DamageDealt,
			DamageTaken:          p.DamageTaken,
			DamageDeltaEnemy:     p.DamageDeltaEnemy,
			DamagePercentageTeam: p.DamagePercentageTeam,
			GoldEarned:           p.GoldEarned,
			GoldDeltaEnemy:       p.GoldDeltaEnemy,
			GoldPercentageTeam:   p.GoldPercentageTeam,
			VisionScore:          p.VisionScore,
			PinkWardsBought:      p.PinkWardsBought,
		},
		Date:       m.Date,
		Duration:   m.Duration,
		Win:        win,
		RankBefore: nil,
		RankAfter:  nil,
	}

	if mapper.BeforeStatus != nil {
		status := PostgresToRankStatus{
			Status: *mapper.BeforeStatus,
			Detail: mapper.BeforeDetail,
		}.Map()
		result.RankBefore = &status
	}

	if mapper.AfterStatus != nil {
		status := PostgresToRankStatus{
			Status: *mapper.AfterStatus,
			Detail: mapper.AfterDetail,
		}.Map()
		result.RankAfter = &status
	}

	return result
}

type PostgresToRankStatus struct {
	Status RankStatus
	Detail *RankDetail
}

func (mapper PostgresToRankStatus) Map() internal.RankStatus {
	status := mapper.Status
	detail := mapper.Detail

	result := internal.RankStatus{
		PUUID:         internal.NewPUUIDFromString(status.PUUID),
		EffectiveDate: status.EffectiveDate,
		Detail:        nil,
	}

	if detail != nil {
		result.Detail = &internal.RankDetail{
			Wins:   detail.Wins,
			Losses: detail.Losses,
			Rank:   internal.Rank{
				Tier:     convertStringToRiotTier(detail.Tier),
				Division: convertStringToRiotRank(detail.Division),
				LP:       detail.LP,
			},
		}
	}

	return result
}

func convertRiotTierToString(tier riot.Tier) string {
	switch tier {
	case riot.TierIron:
		return "Iron"
	case riot.TierBronze:
		return "Bronze"
	case riot.TierSilver:
		return "Silver"
	case riot.TierGold:
		return "Gold"
	case riot.TierPlatinum:
		return "Platinum"
	case riot.TierEmerald:
		return "Emerald"
	case riot.TierDiamond:
		return "Diamond"
	case riot.TierMaster:
		return "Master"
	case riot.TierGrandmaster:
		return "Grandmaster"
	case riot.TierChallenger:
		return "Challenger"
	default:
		panic("bro.")
	}
}

func convertStringToRiotTier(tier string) riot.Tier {
	switch tier {
	default:
		panic("bro.")
	case "Iron":
		return riot.TierIron
	case "Bronze":
		return riot.TierBronze
	case "Silver":
		return riot.TierSilver
	case "Gold":
		return riot.TierGold
	case "Platinum":
		return riot.TierPlatinum
	case "Emerald":
		return riot.TierEmerald
	case "Diamond":
		return riot.TierDiamond
	case "Master":
		return riot.TierMaster
	case "Grandmaster":
		return riot.TierGrandmaster
	case "Challenger":
		return riot.TierChallenger
	}
}

func convertStringToRiotRank(rank string) riot.Division {
	switch rank {
	default:
		panic("bro.")
	case "I":
		return riot.Division1
	case "II":
		return riot.Division2
	case "III":
		return riot.Division3
	case "IV":
		return riot.Division4
	}
}

func convertListToRunePage(runes [11]int) internal.RunePage {
	page := internal.RunePage{
		PrimaryTree:     runes[0],
		PrimaryKeystone: runes[1],
		PrimaryA:        runes[2],
		PrimaryB:        runes[3],
		PrimaryC:        runes[4],
		SecondaryTree:   runes[5],
		SecondaryA:      runes[6],
		SecondaryB:      runes[7],
		MiniOffense:     runes[8],
		MiniFlex:        runes[9],
		MiniDefense:     runes[10],
	}

	return page
}

func convertRunePageToList(runes internal.RunePage) [11]int {
	ids := [11]int{
		runes.PrimaryTree,
		runes.PrimaryKeystone,
		runes.PrimaryA,
		runes.PrimaryB,
		runes.PrimaryC,
		runes.SecondaryTree,
		runes.SecondaryA,
		runes.SecondaryB,
		runes.MiniOffense,
		runes.MiniFlex,
		runes.MiniDefense,
	}

	return ids
}

func convertTeamPositionToString(position internal.TeamPosition) string {
	switch position {
	case internal.TeamPositionBottom:
		return "Bottom"
	case internal.TeamPositionSupport:
		return "Support"
	case internal.TeamPositionTop:
		return "Top"
	case internal.TeamPositionMiddle:
		return "Middle"
	case internal.TeamPositionJungle:
		return "Jungle"
	}

	return ""
}

func convertStringToTeamPosition(s string) internal.TeamPosition {
	switch s {
	case "Bottom":
		return internal.TeamPositionBottom
	case "Support":
		return internal.TeamPositionSupport
	case "Middle":
		return internal.TeamPositionMiddle
	case "Top":
		return internal.TeamPositionTop
	case "Jungle":
		return internal.TeamPositionJungle
	}

	return internal.TeamPositionTop
}
