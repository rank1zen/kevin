package internal

import (
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type MatchID string

func (id MatchID) String() string {
	return string(id)
}

type GameVersion string

func (v GameVersion) GetPatch() string {
	return string(v[0:5])
}

type ParticipantID int

type Team int

const (
	BlueSide Team = 100
	RedSide  Team = 200
)

type Match struct {
	ID       string
	Date     time.Time
	Duration time.Duration
	Version  GameVersion
	Winner   Team
}

func NewMatch(opts ...func(*Match)) Match {
	var match Match
	for _, f := range opts {
		f(&match)
	}
	return match
}

func WithRiotMatch(match *riot.Match) func(*Match) {
	var winner Team
	if match.Info.Teams[0].Win {
		winner = Team(match.Info.Teams[0].TeamId)
	} else {
		winner = Team(match.Info.Teams[1].TeamId)
	}
	return func(m *Match) {
		m.ID = match.Metadata.MatchId
		m.Date = makeRiotUnixTimeStamp(match.Info.GameEndTimestamp)
		m.Duration = makeRiotTimeDuration(match.Info.GameDuration)
		m.Version = GameVersion(match.Info.GameVersion)
		m.Winner = winner
	}
}

type MatchSummoner struct {
	Name string
	Rank *RankDetail

	Participant
}

func makeRiotUnixTimeStamp(ts int64) time.Time {
	return time.Unix(ts/1000, 0)
}

func makeRiotTimeDuration(t int64) time.Duration {
	return time.Second * time.Duration(t)
}
