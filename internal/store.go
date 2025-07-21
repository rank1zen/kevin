package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

var (
	// ErrUnknownStoreError is a generic store error.
	ErrUnknownStoreError = errors.New("unknown store error")

	// ErrSummonerNotFound is returned by Store.GetSummoner when a summoner
	// is not found in the store.
	ErrSummonerNotFound = errors.New("summoner not found")

	// ErrRankUnavailable indicates that the store does not have an
	// available rank record for a summoner. Returned by GetRank.
	ErrRankUnavailable = errors.New("rank unavailable")

	// ErrMatchNotFound is returned by Store.GetMatch when a match is not
	// found in store.
	ErrMatchNotFound = errors.New("rank unavailable")
)

// Store manages persistent data.
type Store interface {
	// GetSummoner returns the summoner, if found in store, otherwise,
	// return ErrSummonerNotFound.
	GetSummoner(ctx context.Context, puuid riot.PUUID) (Summoner, error)

	// GetPUUID returns the summoner's puuid, if found in store, otherwise,
	// return ErrSummonerNotFound.
	GetPUUID(ctx context.Context, name, tag string) (riot.PUUID, error)

	// GetMatch returns the match, if found in store, otherwise,
	// return ...
	GetMatch(ctx context.Context, id riot.PUUID) (Match, error)

	// GetRank returns the most recent rank for a summoner before or at
	// time ts if recent is true, otherwise returns the oldest rank after
	// time ts. Returns ErrRankUnavailable if a rank record satisfying the
	// criteria does not exist.
	GetRank(ctx context.Context, puuid riot.PUUID, ts time.Time, recent bool) (RankRecord, error)

	// GetZMatches returns matches a summoner has played in the given time
	// range.
	GetZMatches(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]SummonerMatch, error)

	// GetChampions returns the summoner champion stats in the given time
	// range.
	GetChampions(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]SummonerChampion, error)

	// GetNewMatchIDs returns the ids of matches not in store.
	GetNewMatchIDs(ctx context.Context, ids []string) (newIDs []string, err error)

	// GetMatches returns a summoners match history in chronological order.
	// Each page is 10 matches.
	//
	// Might be deprecated
	GetMatches(ctx context.Context, puuid string, page int) ([]SummonerMatch, error)

	// RecordMatch creates the match record.
	RecordMatch(ctx context.Context, match Match) error

	// RecordMatchTimeline creates the match timeline event records.
	RecordTimeline(ctx context.Context, id string, items []ItemEvent, skills []SkillEvent) error

	// RecordSummoner updates the summoner and their rank.
	// The rank is set as the most recent available record of rank.
	RecordSummoner(ctx context.Context, summoner Summoner, rank RankStatus) error

	// SearchSummoner returns the best matches for query q.
	SearchSummoner(ctx context.Context, q string) ([]SearchResult, error)
}

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
			m.Participants[i] = NewParticipant(RiotMatchToParticipant(*match, riot.PUUID(p.PUUID)))
		}

		return nil
	}
}

func NewPUUIDFromString(s string) riot.PUUID {
	if len(s) != 78 {
		panic("puuid string not of len 78")
	}

	return riot.PUUID(s)
}

// Participant represents a record of a summoner in a ranked match.
type Participant struct {
	PUUID                  riot.PUUID
	MatchID                string
	TeamID                 int
	ChampionID             int
	ChampionLevel          int
	SummonerIDs            [2]int
	Runes                  RunePage
	Items                  [7]int
	Kills, Deaths, Assists int
	KillParticipation      float32
	CreepScore             int
	CreepScorePerMinute    float32
	DamageDealt            int
	DamageTaken            int
	DamageDeltaEnemy       int
	DamagePercentageTeam   float32
	GoldEarned             int
	GoldDeltaEnemy         int
	GoldPercentageTeam     float32
	VisionScore            int
	PinkWardsBought        int
}

type ParticipantOption func(*Participant) error

func NewParticipant(opts ...ParticipantOption) Participant {
	var p Participant
	for _, f := range opts {
		f(&p)
	}

	return p
}

func RiotMatchToParticipant(match riot.Match, puuid riot.PUUID) ParticipantOption {
	var s *riot.MatchParticipant
	for _, p := range match.Info.Participants {
		if p.PUUID == string(puuid) {
			s = p
		}
	}

	if s == nil {
		panic("yo puuid is not in this match")
	}

	var teamDamage, teamGold, teamKills int
	for _, p := range match.Info.Participants {
		if p.TeamID == s.TeamID {
			teamDamage += p.TotalDamageDealtToChampions
			teamGold += p.GoldEarned
			teamKills += p.Kills
		}
	}

	var counterpart *riot.MatchParticipant
	for _, p := range match.Info.Participants {
		if p.TeamPosition == s.TeamPosition && p.PUUID != s.PUUID {
			counterpart = p
		}
	}

	var (
		cs       = s.TotalMinionsKilled + s.NeutralMinionsKilled
		csPerMin = float32(cs*60) / float32(match.Info.GameDuration)
		kp       = float32(s.Assists+s.Kills) / float32(teamKills)

		damageDelta = s.TotalDamageDealtToChampions - counterpart.TotalDamageDealtToChampions
		goldDelta   = s.GoldEarned - counterpart.GoldEarned
		damageShare = float32(s.TotalDamageDealtToChampions) / float32(teamDamage)
		goldShare   = float32(s.GoldEarned) / float32(teamGold)

		runes     = NewRunePage(WithRiotPerks(s.Perks))
		items     = [7]int{s.Item0, s.Item1, s.Item2, s.Item3, s.Item4, s.Item5, s.Item6}
		summoners = [2]int{s.Summoner1ID, s.Summoner2ID}
	)

	return func(p *Participant) error {
		p.PUUID = NewPUUIDFromString(s.PUUID)
		p.MatchID = match.Metadata.MatchID
		p.TeamID = s.TeamID
		p.ChampionID = s.ChampionID
		p.ChampionLevel = s.ChampLevel
		p.SummonerIDs = summoners
		p.Runes = runes
		p.Items = items
		p.Kills = s.Kills
		p.Deaths = s.Deaths
		p.Assists = s.Assists
		p.KillParticipation = kp
		p.CreepScore = cs
		p.CreepScorePerMinute = csPerMin
		p.DamageDealt = s.TotalDamageDealtToChampions
		p.DamageTaken = s.TotalDamageTaken
		p.DamageDeltaEnemy = damageDelta
		p.DamagePercentageTeam = damageShare
		p.GoldEarned = s.GoldEarned
		p.GoldDeltaEnemy = goldDelta
		p.GoldPercentageTeam = goldShare
		p.VisionScore = s.VisionScore
		p.PinkWardsBought = s.DetectorWardsPlaced
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

type LiveParticipant struct {
	PUUID        string
	MatchID      string
	ChampionID   int
	Runes        RunePage
	TeamID       int
	SummonersIDs [2]int
}

func NewLiveParticipant(opts ...LiveParticipantOption) LiveParticipant {
	var m LiveParticipant
	for _, f := range opts {
		f(&m)
	}
	return m
}

type LiveParticipantOption func(*LiveParticipant) error

func WithRiotCurrentGame(r riot.LiveMatch, puuid string) LiveParticipantOption {
	var selected *riot.LiveMatchParticipant
	for _, p := range r.Participants {
		if p.PUUID == puuid {
			selected = &p
		}
	}
	if selected == nil {
		panic("bro.")
	}

	matchID := fmt.Sprintf("%s_%d", r.PlatformID, r.GameID)

	return func(m *LiveParticipant) error {
		m.PUUID = selected.PUUID
		m.MatchID = matchID
		m.ChampionID = selected.ChampionID
		m.Runes = NewRunePage(WithRiotSpectatorPerks(&selected.Perks))
		m.TeamID = selected.TeamID
		m.SummonersIDs = [2]int{selected.Spell1ID, selected.Spell2ID}
		return nil
	}
}

type Summoner struct {
	PUUID         riot.PUUID
	Name, Tagline string
	Platform      string
	SummonerID    string
}

type SummonerMatch struct {
	Date     time.Time
	Duration time.Duration
	LpDelta  *int
	Win      bool

	Participant
}

// SummonerChampion is a summoner's champion stats averaged over GamesPlayed.
type SummonerChampion struct {
	PUUID riot.PUUID

	// NOTE: Champion type should be specified by ddragon package.
	Champion Champion

	GamesPlayed int

	Wins, Losses int

	AverageKillsPerGame float32

	AverageDeathsPerGame float32

	AverageAssistsPerGame float32

	AverageKillParticipationPerGame float32

	AverageCreepScorePerGame float32

	AverageCreepScorePerMinutePerGame float32

	AverageDamageDealtPerGame float32

	AverageDamageTakenPerGame float32

	AverageDamageDeltaEnemyPerGame float32

	AverageDamagePercentagePerGame float32

	AverageGoldEarnedPerGame float32

	AverageGoldDeltaEnemyPerGame float32

	AverageGoldPercentagePerGame float32

	AverageVisionScorePerGame float32

	AveragePinkWardsBoughtPerGame float32
}

type SearchResult struct {
	Page    string
	Puuid   riot.PUUID
	Name    string
	Tagline string
}

type Rank struct {
	Tier     riot.Tier
	Division riot.Division
	LP       int
}

type RankStatus struct {
	PUUID         riot.PUUID
	EffectiveDate time.Time
	Detail        *RankDetail
}

type RankRecord struct {
	PUUID         riot.PUUID
	EffectiveDate time.Time
	EndDate       *time.Time
	IsCurrent     bool
	Detail        *RankDetail
}

type RankDetail struct {
	Wins     int
	Losses   int
	Tier     riot.Tier
	Division riot.Division
	LP       int
}

func NewRankDetail(opts ...RankDetailOption) RankDetail {
	var m RankDetail
	for _, f := range opts {
		if err := f(&m); err != nil {
			panic(err)
		}
	}
	return m
}

type RankDetailOption func(*RankDetail) error

func WithRiotLeagueEntry(rank riot.LeagueEntry) RankDetailOption {
	return func(m *RankDetail) error {
		m.Wins = rank.Wins
		m.Losses = rank.Losses
		m.Tier = rank.Tier
		m.Division = rank.Division
		m.LP = rank.LeaguePoints
		return nil
	}
}
