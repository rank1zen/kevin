package frontend

import (
	"errors"
	"strings"

	"github.com/rank1zen/kevin/internal/riot"
)

var stringToRiotRegion = map[string]riot.Region{
	"NA1": riot.RegionNA1,
	"KR": riot.RegionKR,
}

func convertStringToRiotRegion(s string) riot.Region {
	if region, found := stringToRiotRegion[s]; found {
		return region
	}

	return riot.RegionNA1
}

var riotRegionToString = map[riot.Region]string{
	riot.RegionNA1: "NA1",
	riot.RegionKR: "KR",
}

func convertRiotRegionToString(region riot.Region) (string) {
	if string, found := riotRegionToString[region]; found {
		return string
	}

	return riotRegionToString[riot.RegionNA1]
}

func convertRiotIDToNameTag(riotID string) (name, tag string, err error) {
	tokens := strings.SplitN(riotID, "-", 2)
	name = tokens[0]
	tag = tokens[1]
	if tag == "" {
		return "", "", errors.New("invalid riotID")
	}

	return name, tag, nil
}
