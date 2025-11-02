package internal

import (
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

func findSoloQLeagueEntry(entries riot.LeagueList) (soloq *riot.LeagueEntry) {
	for _, entry := range entries {
		if entry.QueueType == riot.QueueTypeRankedSolo5x5 {
			return &entry
		}
	}

	return nil
}

func soloQMatchFilter(start, end time.Time) riot.MatchListOptions {
	options := riot.MatchListOptions{
		StartTime: new(int64),
		EndTime:   new(int64),
		Queue:     new(int),
		Type:      nil,
		Start:     0,
		Count:     100,
	}

	*options.Queue = 420
	*options.StartTime = start.Unix()
	*options.EndTime = end.Unix()

	return options
}
