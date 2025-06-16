package internal

import (
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type Match struct {
	ID       string

	// Date is the end timestamp
	Date     time.Time

	Duration time.Duration

	Version  string

	WinnerID int
}

func NewMatch(opts ...MatchOption) Match {
	var match Match
	for _, f := range opts {
		f(&match)
	}
	return match
}

type MatchOption func(*Match) error

func WithRiotMatch(match *riot.Match) MatchOption {
	var winner int
	if match.Info.Teams[0].Win {
		winner = match.Info.Teams[0].TeamId
	} else {
		winner = match.Info.Teams[1].TeamId
	}
	return func(m *Match) error {
		m.ID = match.Metadata.MatchId
		m.Date = makeRiotUnixTimeStamp(match.Info.GameEndTimestamp)
		m.Duration = makeRiotTimeDuration(match.Info.GameDuration)
		m.Version = match.Info.GameVersion
		m.WinnerID = winner
		return nil
	}
}

type LiveMatch struct {
	ID       string

	// Date is game start timestamp
	Date     time.Time

	Participants [10]LiveParticipant
}

type MatchSummoner struct {
	Name string
	Rank *RankDetail

	Participant
}

func makeRiotUnixTimeStamp(ts int64) time.Time {
	return time.UnixMilli(ts)
}

func makeRiotTimeDuration(t int64) time.Duration {
	return time.Second * time.Duration(t)
}
