package service

import (
	"context"
	"errors"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

type LiveMatchService Service

type GetLiveMatchRequest struct {
	Region *riot.Region `json:"region"`
	PUUID  riot.PUUID   `json:"puuid"`
}

func (s *LiveMatchService) GetLiveMatch(ctx context.Context, req GetLiveMatchRequest) (*internal.LiveMatch, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
	}

	riotGame, err := s.riot.Spectator.GetLiveMatch(ctx, *req.Region, req.PUUID.String())
	if err != nil {
		if errors.Is(err, riot.ErrNotFound) {
			return nil, ErrNoLiveMatch
		}

		return nil, err
	}

	match := internal.RiotToLiveMatchMapper{Match: *riotGame}.Map()

	return &match, nil
}

type GetLiveMatchByIDRequest struct {
	Region  *riot.Region `json:"region"`
	MatchID string       `json:"match_id"`
}

func (s *LiveMatchService) GetLiveMatchByID(ctx context.Context, req GetLiveMatchByIDRequest) (*LiveMatchDetail, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
	}

	storeStatus, err := s.store.LiveMatchStatus.GetLiveMatchStatus(ctx, *req.Region, req.MatchID)
	if err != nil {
		return nil, err
	}

	storeLiveMatch, err := s.store.LiveMatch.GetLiveMatch(ctx, req.MatchID)
	if err != nil {
		if errors.Is(err, ErrNoLiveMatch) {
			return nil, ErrNoLiveMatch
		}

		return nil, err
	}

	result := LiveMatchDetail{
		ID:           storeLiveMatch.ID,
		Date:         storeLiveMatch.Date,
		Participants: [10]LiveParticipantDetail{},
		Expired:      storeStatus.Expired,
	}

	for i := range result.Participants {
		storeParticipant := storeLiveMatch.Participants[i]
		result.Participants[i] = LiveParticipantDetail{
			PUUID:        storeParticipant.PUUID,
			MatchID:      storeLiveMatch.ID,
			ChampionID:   storeParticipant.ChampionID,
			Runes:        storeParticipant.Runes,
			SummonersIDs: storeParticipant.SummonersIDs,
			TeamID:       storeParticipant.TeamID,
			CurrentRank:  nil, // TODO: Implement rank retrieval logic
		}
	}

	return &result, nil
}

// LiveMatchDetail represents details of a live match.
type LiveMatchDetail struct {
	Region riot.Region

	ID string

	Date time.Time

	Participants [10]LiveParticipantDetail

	// Expired indicates whether the match has ended.
	Expired bool
}

// LiveParticipantDetail represents a participant in a live match.
type LiveParticipantDetail struct {
	PUUID riot.PUUID

	MatchID string

	ChampionID int

	Runes internal.RunePage

	SummonersIDs [2]int

	TeamID int

	CurrentRank *internal.Rank
}
