package postgres

import (
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
)

type PostgresToRankStatus struct {
	Status postgres.RankStatus
	Detail *postgres.RankDetail
}

func (mapper PostgresToRankStatus) Map() internal.RankStatus {
	status := mapper.Status
	detail := mapper.Detail

	result := internal.RankStatus{
		PUUID:         riot.PUUID(status.PUUID),
		EffectiveDate: status.EffectiveDate,
		Detail:        nil,
	}

	if detail != nil {
		result.Detail = &internal.RankDetail{
			Wins:   detail.Wins,
			Losses: detail.Losses,
			Rank: internal.Rank{
				Tier:     convertStringToRiotTier(detail.Tier),
				Division: convertStringToRiotRank(detail.Division),
				LP:       detail.LP,
			},
		}
	}

	return result
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
