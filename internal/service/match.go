package service

import (
	"context"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

// MatchService manages match-related operations.
type MatchService Service

// GetMatchlistRequest represents the request payload for retrieving a list of matches.
type GetMatchlistRequest struct {
	Region  *riot.Region `json:"region"`
	PUUID   riot.PUUID   `json:"puuid"`
	StartTS *time.Time   `json:"startTs"`
	EndTS   *time.Time   `json:"endTs"`
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

		match := mapRiotMatchToModelMatch(*riotMatch)

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

	newIDS, err := s.match.GetNewMatchIDs(ctx, []string{req.MatchID})
	if err != nil {
		return nil, err
	}

	if len(newIDS) == 1 { // If the match is new, fetch it from Riot API and record it
		riotMatch, err := s.riot.Match.GetMatch(ctx, *req.Region, req.MatchID)
		if err != nil {
			return nil, err
		}

		match := mapRiotMatchToModelMatch(*riotMatch)

		err = s.match.RecordMatch(ctx, match)
		if err != nil {
			return nil, err
		}
	}

	storeMatch, err := s.match.GetMatchDetail(ctx, req.MatchID)
	if err != nil {
		return nil, err
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

// mapRiotMatchToModelMatch converts a riot.Match struct to a model.Match struct.
// This is a placeholder for the actual mapping logic, which might involve a dedicated mapper package.
func mapRiotMatchToModelMatch(riotMatch riot.Match) internal.Match {
	// A highly simplified placeholder for demonstration.
	// Real implementation would involve mapping all relevant fields and nested structs.
	return internal.Match{
		ID:       riotMatch.Metadata.MatchID,
		Date:     time.Unix(0, riotMatch.Info.GameEndTimestamp*int64(time.Millisecond)),
		Duration: time.Duration(riotMatch.Info.GameDuration) * time.Second,
		Version:  riotMatch.Info.GameVersion,
		// WinnerID:      ... (requires mapping logic)
		// Participants:  ... (requires mapping logic for each participant)
	}
}
