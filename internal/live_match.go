package internal

import (
	"context"
	"errors"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type LiveMatch struct {
	ID string

	// Date is game start timestamp
	Date time.Time

	// Participants are the players in this current match. There is no
	// chosen order.
	Participants [10]LiveParticipant
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

type LiveMatchService Datasource

type GetLiveMatchRequest struct {
	Region *riot.Region `json:"region"`
	PUUID  riot.PUUID   `json:"puuid"`
}

func (s *LiveMatchService) GetLiveMatch(ctx context.Context, req GetLiveMatchRequest) (*LiveMatch, error) {
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

	match := RiotToLiveMatchMapper{Match: *riotGame}.Map()

	return &match, nil
}
