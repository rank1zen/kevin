package internal

import (
	"fmt"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

// Match represents a record of a ranked match.
type Match struct {
	// ID is a region+number, which forms an identifier.
	ID string

	// Date is the end timestamp of the match.
	Date time.Time

	// Duration is the length of the match.
	Duration time.Duration

	// Version is the game version.
	Version string

	// WinnerID is the ID of the winning team.
	WinnerID int

	Participants [10]Participant
}

func NewMatch(opts ...MatchOption) Match {
	var match Match
	for _, f := range opts {
		if err := f(&match); err != nil {
			panic(err)
		}
	}

	return match
}

type MatchOption func(*Match) error

func WithRiotMatch(match *riot.Match) MatchOption {
	return func(m *Match) error {
		m.ID = match.Metadata.MatchID
		m.Date = makeRiotUnixTimeStamp(match.Info.GameEndTimestamp)
		m.Duration = makeRiotTimeDuration(match.Info.GameDuration)
		m.Version = match.Info.GameVersion

		var winner int
		if match.Info.Teams[0].Win {
			winner = match.Info.Teams[0].TeamID
		} else {
			winner = match.Info.Teams[1].TeamID
		}

		m.WinnerID = winner

		for i, p := range match.Info.Participants {
			m.Participants[i] = NewParticipant(RiotMatchToParticipant(*match, p.PUUID))
		}

		return nil
	}
}

type LiveMatch struct {
	ID string

	// Date is game start timestamp
	Date time.Time

	Participants [10]LiveParticipant
}

func NewLiveMatch(opts ...LiveMatchOption) LiveMatch {
	var match LiveMatch

	for _, f := range opts {
		f(&match)
	}

	return match
}

type LiveMatchOption func(*LiveMatch) error

func WithRiotLiveMatch(match *riot.LiveMatch) LiveMatchOption {
	matchID := fmt.Sprintf("%s_%d", match.PlatformID, match.GameID)

	participants := []LiveParticipant{}
	for _, p := range match.Participants {
		participants = append(participants, LiveParticipant{
			PUUID:        p.PUUID,
			MatchID:      matchID,
			ChampionID:   p.ChampionID,
			Runes:        NewRunePage(WithRiotSpectatorPerks(&p.Perks)),
			TeamID:       p.TeamID,
			SummonersIDs: [2]int{p.Spell1ID, p.Spell2ID},
		})
	}

	return func(m *LiveMatch) error {
		m.ID = matchID
		m.Date = makeRiotUnixTimeStamp(match.GameStartTime)
		m.Participants = [10]LiveParticipant(participants)
		return nil
	}
}

func makeRiotUnixTimeStamp(ts int64) time.Time {
	return time.UnixMilli(ts)
}

func makeRiotTimeDuration(t int64) time.Duration {
	return time.Second * time.Duration(t)
}
