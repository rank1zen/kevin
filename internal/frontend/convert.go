package frontend

import (
	"github.com/rank1zen/kevin/internal/riot"
)

var stringToRiotRegion = map[string]riot.Region{
	"NA1": riot.RegionNA1,
	"KR":  riot.RegionKR,
}

func convertStringToRiotRegion(s string) riot.Region {
	if region, found := stringToRiotRegion[s]; found {
		return region
	}

	return riot.RegionNA1
}
