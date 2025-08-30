package internal

import (
	"time"
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
