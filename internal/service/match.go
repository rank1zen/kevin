package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

// MatchService manages match-related operations.
type MatchService Service

// GetMatchlistRequest represents the request payload for retrieving a list of matches.
type GetMatchlistRequest struct {
	Region *riot.Region `json:"region"`

	PUUID riot.PUUID `json:"puuid"`

	// StartTS is the start timestamp of the end of the game from which to include
	// in the match list. Defaults to 1 day ago.
	StartTS *time.Time `json:"startTs"`

	// EndTS is the end timestamp of the end of the game from which to include in
	// the match list. Defaults to now.
	EndTS *time.Time `json:"endTs"`
}

// GetMatchlist retrieves a list of matches for a given summoner within a time range.
func (s *MatchService) GetMatchlist(ctx context.Context, req GetMatchlistRequest) ([]internal.SummonerMatch, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1 // Default to NA1 if not specified
	}

	currTime := time.Now().In(time.UTC)

	if req.StartTS == nil {
		req.StartTS = new(time.Time)
		*req.StartTS = currTime.AddDate(0, 0, -1) // Default to 1 day ago
	}

	if req.EndTS == nil {
		req.EndTS = new(time.Time)
		*req.EndTS = currTime // Default to current time
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

	newIDs, err := s.store.Match.GetNewMatchIDs(ctx, matchIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get new match IDs in store: %w", err)
	}

	// TODO: put these in batch
	for _, id := range newIDs {
		riotMatch, err := s.riot.Match.GetMatch(ctx, *req.Region, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get match details from Riot: %w", err)
		}

		match := internal.RiotToMatchMapper{Match: *riotMatch}.Map()

		err = s.store.Match.RecordMatch(ctx, match)
		if err != nil {
			return nil, fmt.Errorf("failed to record match in store: %w", err)
		}
	}

	storeMatches, err := s.store.Match.GetMatchlist(ctx, req.PUUID, *req.StartTS, *req.EndTS)
	if err != nil {
		return nil, fmt.Errorf("failed to get matchlist from store: %w", err)
	}

	return storeMatches, nil
}

// GetMatchDetailRequest represents the request payload for retrieving a single match's details.
type GetMatchDetailRequest struct {
	Region  *riot.Region `json:"region"`
	MatchID string       `json:"matchId"`
}

// GetMatchDetail retrieves the details of a specific match.
func (s *MatchService) GetMatchDetail(ctx context.Context, req GetMatchDetailRequest) (*internal.MatchDetail, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1 // Default to NA1 if not specified
	}

	newIDS, err := s.store.Match.GetNewMatchIDs(ctx, []string{req.MatchID})
	if err != nil {
		return nil, fmt.Errorf("failed to check match IDs in store: %w", err)
	}

	if len(newIDS) == 1 {
		// If the match is new, fetch it from Riot API and record it
		riotMatch, err := s.riot.Match.GetMatch(ctx, *req.Region, req.MatchID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch match from Riot API: %w", err)
		}

		match := internal.RiotToMatchMapper{Match: *riotMatch}.Map()

		err = s.store.Match.RecordMatch(ctx, match)
		if err != nil {
			return nil, fmt.Errorf("failed to record match in store: %w", err)
		}
	}

	storeMatch, err := s.store.Match.GetMatchDetail(ctx, req.MatchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match detail from store: %w", err)
	}

	return &storeMatch, nil
}

func soloQMatchFilter(start, end time.Time) riot.MatchListOptions {
	options := riot.MatchListOptions{
		StartTime: new(int64),
		EndTime:   new(int64),
		Queue:     new(int),
		Type:      nil,
		Start:     0,
		Count:     100,
	}

	*options.Queue = 420
	*options.StartTime = start.Unix()
	*options.EndTime = end.Unix()

	return options
}
