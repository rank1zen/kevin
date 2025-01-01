package internal

import "time"

// LiveParticipant is a domain type representing a live match participant. Many
// fields from Participant is not yet availble.
type LiveParticipant struct {
	Puuid          PUUID
	TeamID         TeamID
	StartTimestamp time.Time
	Duration       time.Duration
	BannedChampion *ChampionID
	Champion       ChampionID
	Summs          SummsIDs
	Runes          Runes
}
