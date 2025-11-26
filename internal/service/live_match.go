package service

import (
	"context"
	"errors"

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
