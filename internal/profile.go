package internal

import (
	"fmt"
	"time"
)

// Profile represents a record of a summoner.
//
// Name and Tagline will be always be the current name and tagline for a
// profile.
type Profile struct {
	ValidFrom  time.Time
	ValidTo    time.Time
	RecordDate time.Time

	Puuid         PUUID
	SummonerID    SummonerID
	AccountID     AccountID
	RevisionDate  time.Time
	Level         int
	ProfileIconID ProfileIconID
	Rank          *RankRecord
	Name          string
	Tagline       string
}

type SummonerRecord struct {
	ValidFrom time.Time
	ValidTo   time.Time
	EnteredAt time.Time

	Puuid         PUUID
	ID            SummonerID
	RevisionDate  time.Time
	Level         int
	ProfileIconID ProfileIconID
}

type RankRecord struct {
	ValidFrom time.Time
	ValidTo   time.Time
	Timestamp time.Time

	Puuid    PUUID
	LeagueID string
	Wins     int
	Losses   int
	Tier     string
	Division string
	LP       int
}

// TODO: implement the cases for Challenger, GM, Masters.
func (r *RankRecord) RankString() string {
	if r == nil {
		return "Unranked"
	}

	return fmt.Sprintf("%s %s %dLP", r.Tier, r.Division, r.LP)
}

// RankSnapshot is a record of a rank prior to a match, and LpDelta is the LP
// change after the match, if available.
type RankSnapshot struct {
	RankRecord

	LpDelta *int
}
