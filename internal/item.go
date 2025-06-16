package internal

import (
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type ItemEventType int

const (
	ItemPurchasedEvent = iota
	ItemSoldEvent
)

// change this to ItemFrame
type ItemEvent struct {
	MatchID         string
	PUUID           int
	ItemID          int
	InGameTimestamp time.Duration
	Type            ItemEventType
}

type ItemFrame struct {
	// InGameTimestamp is the time of the first item event
	InGameTimestamp time.Duration

	// Events should be at most 30 seconds apart
	Events []ItemEvent
}

// ItemProgression is a summoner's item build order
type ItemProgression []ItemFrame

func makeItemProgression(events []ItemEvent) ItemProgression {
	items := []ItemFrame{}
	if len(events) == 0 {
		return items
	}

	newBucket := func(i int) {
		items = append(
			items,
			ItemFrame{
				InGameTimestamp: events[i].InGameTimestamp,
				Events:          []ItemEvent{events[i]},
			},
		)
	}

	newBucket(0)
	curr := events[0].InGameTimestamp
	bucket := items[0]
	j := 0
	for i := 1; i < len(events); i++ {
		if events[i].InGameTimestamp-curr < 30*time.Second {
			curr = events[i].InGameTimestamp
			bucket.Events = append(bucket.Events, events[i])
		} else {
			newBucket(i)
			j++
			bucket = items[j]
		}
	}
	return items
}

func makeItemEvents(timeline riot.Timeline) []ItemEvent {
	panic("not imlpemented")
}
