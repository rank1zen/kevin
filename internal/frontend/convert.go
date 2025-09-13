package frontend

import (
	"github.com/rank1zen/kevin/internal/riot"
)

var stringToRiotRegion = map[string]riot.Region{
	"BR1":  riot.RegionBR1,
	"EUN1": riot.RegionEUN1,
	"EUW1": riot.RegionEUW1,
	"JP1":  riot.RegionJP1,
	"KR":   riot.RegionKR,
	"LA1":  riot.RegionLA1,
	"LA2":  riot.RegionLA2,
	"ME1":  riot.RegionME1,
	"NA1":  riot.RegionNA1,
	"OC1":  riot.RegionOC1,
	"RU":   riot.RegionRU,
	"SG2":  riot.RegionSEA,
	"TR1":  riot.RegionTR1,
	"TW2":  riot.RegionTW2,
	"VN2":  riot.RegionVN2,
}

func convertStringToRiotRegion(s string) riot.Region {
	if region, found := stringToRiotRegion[s]; found {
		return region
	}

	return riot.RegionNA1
}
