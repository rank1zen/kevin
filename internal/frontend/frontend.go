package frontend

import (
	"embed"
	"strings"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

//go:embed static
var StaticAssets embed.FS

// GetDays returns 8 datetimes, representing the past 7 days, including the
// given timestamp. Each date time is at midnight in the provided timezone. The
// slice is ordered in most recent order.
func GetDays(ts time.Time) []time.Time {
	year, month, day := ts.Date()
	ref := time.Date(year, month, day, 0, 0, 0, 0, ts.Location())

	days := []time.Time{}

	for i := range 8 {
		offset := -i + 1
		days = append(days, ref.AddDate(0, 0, offset))
	}

	return days
}

func ParseRiotID(riotID string) (name, tag string) {
	index := strings.Index(riotID, "-")
	if index == -1 {
		return riotID, ""
	}

	name = riotID[:index]
	tag = riotID[index+1:]

	return name, tag
}

func StrToRiotRegion(s string) riot.Region {
	if region, found := stringToRiotRegion[s]; found {
		return region
	}

	return riot.RegionNA1
}

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
