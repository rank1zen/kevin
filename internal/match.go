package internal

import (
	"context"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type MatchStore interface {
	RecordMatch(ctx context.Context, match Match) error

	GetMatchlist(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]SummonerMatch, error)

	GetMatchDetail(ctx context.Context, id string) (MatchDetail, error)

	// GetNewMatchIDs returns the ids of matches not in store.
	GetNewMatchIDs(ctx context.Context, ids []string) (newIDs []string, err error)

	GetChampions(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]SummonerChampion, error)
}

// Match represents a record of a ranked match.
type Match struct {
	ID           string
	Date         time.Time
	Duration     time.Duration
	Version      string
	WinnerID     int
	Participants [10]Participant
}

// Participant represents a record of a summoner in a ranked match.
type Participant struct {
	PUUID                  riot.PUUID
	MatchID                string
	TeamID                 int
	ChampionID             int
	ChampionLevel          int
	TeamPosition           TeamPosition
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

// ParticipantDetail is the details relating to a participant record.
type ParticipantDetail struct {
	Participant

	Name, Tag   string
	CurrentRank *RankStatus
	RankBefore  *RankStatus
	RankAfter   *RankStatus
}

type SummonerMatch struct {
	Participant

	Date     time.Time
	Duration time.Duration
	Win      bool

	// RankBefore is the summoner's rank just before the match. A nil value
	// indicates this no record was taken.
	RankBefore *RankStatus
	// RankBefore is the summoner's rank just after the match. A nil value
	// indicates this no record was taken.
	RankAfter *RankStatus
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
