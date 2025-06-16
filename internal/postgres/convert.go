package postgres

import (
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

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
