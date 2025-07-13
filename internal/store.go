package internal

import (
	"context"
	"errors"
	"time"
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
	// GetPUUID returns a summoner's puuid.
	GetPUUID(ctx context.Context, name, tag string) (puuid string, err error)

	GetMatch(ctx context.Context, id string) (Match, [10]Participant, error)

	// GetRank returns the most recent rank for a summoner before or at
	// time ts if recent is true, otherwise returns the oldest rank after
	// time ts. Returns ErrRankUnavailable if a rank record satisfying the
	// criteria does not exist.
	GetRank(ctx context.Context, puuid string, ts time.Time, recent bool) (RankRecord, error)

	// ... returns ErrSummonerNotFound when summoner is not in store.
	GetSummoner(ctx context.Context, puuid string) (Summoner, error)

	// GetMatches returns a summoners match history in chronological order.
	// Each page is 10 matches.
	//
	// Might be deprecated
	GetMatches(ctx context.Context, puuid string, page int) ([]SummonerMatch, error)

	// GetZMatches returns matches a summoner has played in the given time
	// range.
	GetZMatches(ctx context.Context, puuid string, start, end time.Time) ([]SummonerMatch, error)

	// GetChampions returns summoner champions stats in the given time range.
	GetChampions(ctx context.Context, puuid string, start, end time.Time) ([]SummonerChampion, error)

	// GetNewMatchIDs returns the ids of matches not in store.
	GetNewMatchIDs(ctx context.Context, ids []string) (newIDs []string, err error)

	// RecordMatch creates the match records.
	RecordMatch(ctx context.Context, match Match) error

	// RecordMatchTimeline creates the match timeline event records.
	RecordTimeline(ctx context.Context, id string, items []ItemEvent, skills []SkillEvent) error

	// RecordSummoner updates the summoner and their rank.
	// The rank is set as the most recent available record of rank.
	RecordSummoner(ctx context.Context, summoner Summoner, rank RankStatus) error

	// SearchSummoner returns the best matches for query q.
	SearchSummoner(ctx context.Context, q string) ([]SearchResult, error)
}
