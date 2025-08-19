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
	// TODO: rename to ErrMatchNotInStore.
	ErrMatchNotFound = errors.New("match not found")

	ErrMatchMissingParticipants = errors.New("match missing participants")
)

// Store manages persistent data.
type Store interface {
	RecordProfile(ctx context.Context, summoner Profile) error

	RecordMatch(ctx context.Context, match Match) error

	GetProfileDetail(ctx context.Context, puuid riot.PUUID) (ProfileDetail, error)

	// GetMatchDetail returns match details. It returns [ErrMatchNotFound]
	// if the match with id is not found in store.
	GetMatchDetail(ctx context.Context, id string) (MatchDetail, error)

	// GetMatchHistory returns matches a summoner has played in the given
	// time range.
	GetMatchHistory(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]SummonerMatch, error)

	// GetChampions returns averaged stats for each champion played in the
	// given time range.
	GetChampions(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]SummonerChampion, error)

	// GetNewMatchIDs returns the ids of matches not in store.
	GetNewMatchIDs(ctx context.Context, ids []string) (newIDs []string, err error)

	SearchSummoner(ctx context.Context, q string) ([]SearchResult, error)
}

// Match represents a record of a ranked match.
type Match struct {
	// ID is a region+number, which forms an identifier. NOTE: should
	// switch to new match ID type.
	ID string

	// Date is the end timestamp of the match.
	Date time.Time

	// Duration is the length of the match.
	Duration time.Duration

	// Version is the game version.
	Version string

	// WinnerID is the ID of the winning team. NOTE: should switch to new
	// TeamID type.
	WinnerID int

	// Participants are the players in this match. There is no chosen
	// order.
	Participants [10]Participant
}

// TODO: behaviour undocumented.
func (m Match) GetTeamParticipants(teamID int) [5]Participant {
	result := [5]Participant{}
	i := 0
	for _, p := range m.Participants {
		if p.TeamID == teamID {
			result[i] = p
			i++
		}
	}

	return result
}

// MatchDetail are details relating to a match record.
type MatchDetail struct {
	// ID is a region+number, which forms an identifier. NOTE: should
	// switch to new match ID type.
	ID string

	// Date is the end timestamp of the match.
	Date time.Time

	// Duration is the length of the match.
	Duration time.Duration

	// Version is the game version.
	Version string

	// WinnerID is the ID of the winning team. NOTE: should switch to new
	// TeamID type.
	WinnerID int

	// Participants are the players in this match. There is no chosen
	// order.
	Participants [10]ParticipantDetail
}

// Participant represents a record of a summoner in a ranked match.
type Participant struct {
	PUUID riot.PUUID

	MatchID string

	TeamID int

	ChampionID int

	ChampionLevel int

	TeamPosition TeamPosition

	SummonerIDs [2]int

	Runes RunePage

	// Items are in standard inventory order.
	Items [7]int

	Kills, Deaths, Assists int

	KillParticipation float32

	CreepScore int

	CreepScorePerMinute float32

	DamageDealt int

	DamageTaken int

	DamageDeltaEnemy int

	DamagePercentageTeam float32

	GoldEarned int

	GoldDeltaEnemy int

	GoldPercentageTeam float32

	VisionScore int

	PinkWardsBought int
}

// ParticipantDetail is the details relating to a participant record.
type ParticipantDetail struct {
	Participant

	// Name and Tag is the current riot id of the summoner.
	Name, Tag string

	// CurrentRank is the current rank of the summoner.
	CurrentRank *RankStatus

	RankBefore *RankStatus

	RankAfter *RankStatus
}

type SummonerMatch struct {
	Participant

	Date time.Time

	Duration time.Duration

	Win bool

	// RankBefore is the summoner's rank just before the match. A nil value
	// indicates this no record was taken.
	RankBefore *RankStatus

	// RankBefore is the summoner's rank just after the match. A nil value
	// indicates this no record was taken.
	RankAfter *RankStatus
}

type LiveMatch struct {
	ID string

	// Date is game start timestamp
	Date time.Time

	// Participants are the players in this current match. There is no
	// chosen order.
	Participants [10]LiveParticipant
}

// GetTeamParticipants returns a team in the order returned by the riot API,
// which is pick order, and is not necessarily standard order: TOP, JNG, MID,
// BOT, SUP.
func (m LiveMatch) GetTeamParticipants(teamID int) [5]LiveParticipant {
	result := [5]LiveParticipant{}
	if teamID == 100 {
		for i := 0; i < 5; i++ {
			result[i] = m.Participants[i]
		}
	} else {
		for i := 5; i < 10; i++ {
			result[i-5] = m.Participants[i]
		}
	}

	return result
}

// LiveParticipant are currently in a match. NOTE: there are not a lot of
// fields are not available in an on-going game, including the summoners
// position.
type LiveParticipant struct {
	PUUID riot.PUUID

	MatchID string

	ChampionID int

	Runes RunePage

	SummonersIDs [2]int

	TeamID int
}

// RankDetail contains details relating to a summoner's rank.
type RankDetail struct {
	Wins, Losses int

	Rank Rank
}

// RankStatus indicates the status of a summoner's rank.
type RankStatus struct {
	PUUID riot.PUUID

	// EffectiveDate indicates the time this status was taken.
	EffectiveDate time.Time

	// Detail is rank detail. A nil value indicates the summoner is
	// unranked.
	Detail *RankDetail
}

type Profile struct {
	PUUID riot.PUUID

	Name, Tagline string

	Rank RankStatus
}

type ProfileDetail struct {
	PUUID riot.PUUID

	Name, Tagline string

	Rank RankStatus
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
	PUUID riot.PUUID

	Name, Tagline string

	// Rank is the summoners most recent rank record in store.
	Rank *RankStatus
}

func NewPUUIDFromString(s string) riot.PUUID {
	if len(s) != 78 {
		panic("puuid string not of len 78")
	}

	return riot.PUUID(s)
}

func NewPUUID(s string) (riot.PUUID, error) {
	if len(s) != 78 {
		return "", errors.New("puuid string not of len 78")
	}

	return riot.PUUID(s), nil
}

var teamPositions = map[string]TeamPosition{
	"TOP":     0,
	"JUNGLE":  1,
	"MIDDLE":  2,
	"BOTTOM":  3,
	"UTILITY": 4,
}

func convertRiotTeamPosition(s string) TeamPosition {
	pos, ok := teamPositions[s]
	if !ok {
		panic(fmt.Sprintf("team position %s is not valid", s))
	}

	return pos
}
