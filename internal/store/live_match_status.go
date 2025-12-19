package postgres

import (
	"context"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
)

type LiveMatchStatusStore Store

func (s *LiveMatchStatusStore) GetLiveMatchStatus(ctx context.Context, region riot.Region, id string) (*internal.LiveMatchStatus, error) {
	objs := postgres.LiveMatchStatusObjects{Tx: s.Pool}

	liveMatch, err := objs.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	result := internal.LiveMatchStatus{
		Region:  liveMatch.Region,
		ID:      liveMatch.ID,
		Date:    liveMatch.Date,
		Expired: liveMatch.Expired,
	}

	return &result, nil
}

func (s *LiveMatchStatusStore) CreateLiveMatchStatus(ctx context.Context, status *internal.LiveMatchStatus) error {
	objs := postgres.LiveMatchStatusObjects{Tx: s.Pool}

	result := postgres.LiveMatchStatus{
		Region:  status.Region,
		ID:      status.ID,
		Date:    status.Date,
		Expired: status.Expired,
	}

	err := objs.Create(ctx, &result)
	if err != nil {
		return err
	}

	return nil
}

func (s *LiveMatchStatusStore) ExpireLiveMatch(ctx context.Context, region riot.Region, id string) error {
	objs := postgres.LiveMatchStatusObjects{Tx: s.Pool}

	update := postgres.LiveMatchStatusUpdate{}
	update.Expired = new(bool)
	*update.Expired = true

	err := objs.Update(ctx, id, update)

	return err
}
