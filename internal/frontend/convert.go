package frontend

import "github.com/rank1zen/kevin/internal/riot"

var stringToRiotRegion = map[string]riot.Region{
	"NA1": riot.RegionNA1,
}

func convertStringToRiotRegion(s string) (riot.Region, error) {
	region, found := stringToRiotRegion[s]
	if found != true {
		return riot.RegionNA1, ErrInvalidRegion
	}

	return region, nil
}
