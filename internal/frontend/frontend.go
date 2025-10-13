package frontend

import (
	"embed"
	"errors"
	"strings"
	"time"
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

func ParseRiotID(riotID string) (name, tag string, err error) {
	index := strings.Index(riotID, "-")
	if index == -1 {
		return "", "", errors.New("invalid riot id")
	}

	if index == len(riotID)-1 {
		return "", "", errors.New("invalid riot id")
	}

	name = riotID[:index]
	tag = riotID[index+1:]

	if index := strings.Index(tag, "-"); index != -1 {
		return "", "", errors.New("invalid riot id")
	}

	return name, tag, nil
}
