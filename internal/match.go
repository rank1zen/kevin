package internal

import (
	"context"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

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

type MatchStore interface {
	RecordMatch(ctx context.Context, match Match) error

	GetMatchlist(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]SummonerMatch, error)

	GetMatchDetail(ctx context.Context, id string) (MatchDetail, error)

	// GetNewMatchIDs returns the ids of matches not in store.
	GetNewMatchIDs(ctx context.Context, ids []string) (newIDs []string, err error)

	GetChampions(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]SummonerChampion, error)
}

type MatchService Datasource

type GetMatchlistRequest struct {
	Region  *riot.Region `json:"region"`
	PUUID   riot.PUUID   `json:"puuid"`
	StartTS *time.Time   `json:"startTs"`
	EndTS   *time.Time   `json:"endTs"`
}

func (s *MatchService) GetMatchlist(ctx context.Context, req GetMatchlistRequest) ([]SummonerMatch, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
	}

	currTime := time.Now().In(time.UTC)

	if req.StartTS == nil {
		req.StartTS = new(time.Time)
		*req.StartTS = currTime.AddDate(0, 0, -1)
	}

	if req.EndTS == nil {
		req.EndTS = new(time.Time)
		*req.EndTS = currTime
	}

	options := soloQMatchFilter(*req.StartTS, *req.EndTS)
	ids, err := s.riot.Match.GetMatchList(ctx, *req.Region, req.PUUID.String(), options)
	if err != nil {
		return nil, err
	}

	matchIDs := []string{}
	for _, id := range ids {
		matchIDs = append(matchIDs, id)
	}

	newIDs, err := s.match.GetNewMatchIDs(ctx, matchIDs)
	if err != nil {
		return nil, err
	}

	// TODO: put these in batch
	for _, id := range newIDs {
		riotMatch, err := s.riot.Match.GetMatch(ctx, *req.Region, id)
		if err != nil {
			return nil, err
		}

		match := RiotToMatchMapper{Match: *riotMatch}.Map()

		err = s.match.RecordMatch(ctx, match)
		if err != nil {
			return nil, err
		}
	}

	storeMatches, err := s.match.GetMatchlist(ctx, req.PUUID, *req.StartTS, *req.EndTS)
	if err != nil {
		return nil, err
	}

	return storeMatches, nil
}

type GetMatchDetailRequest struct {
	Region  *riot.Region `json:"region"`
	MatchID string       `json:"matchId"`
}

func (s *MatchService) GetMatchDetail(ctx context.Context, req GetMatchDetailRequest) (*MatchDetail, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
	}

	storeMatch, err := s.match.GetMatchDetail(ctx, req.MatchID)
	if err != nil {
		return nil, err
	}

	return &storeMatch, nil
}

type GetSummonerChampionsRequest struct {
	Region  *riot.Region `json:"region"`
	PUUID   riot.PUUID   `json:"puuid"`
	StartTS *time.Time   `json:"startTs"`
	EndTS   *time.Time   `json:"endTs"`
}

func (s *ProfileService) GetSummonerChampions(ctx context.Context, req GetSummonerChampionsRequest) ([]SummonerChampion, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
	}

	currTime := time.Now().In(time.UTC)

	if req.StartTS == nil {
		req.StartTS = new(time.Time)
		*req.StartTS = currTime.AddDate(0, 0, -7)
	}

	if req.EndTS == nil {
		req.EndTS = new(time.Time)
		*req.EndTS = currTime
	}

	storeChamps, err := s.match.GetChampions(ctx, req.PUUID, *req.StartTS, *req.EndTS)
	if err != nil {
		return nil, err
	}

	return storeChamps, nil
}
