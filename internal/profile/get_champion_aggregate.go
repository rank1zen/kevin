package profile

import "context"

type GetChampionAggregateRequest struct {
	Region  string
	MatchID string
}

type GetChampionAggregateResponse struct {
}

func (s *Service) GetChampionAggregate(ctx context.Context, req GetChampionAggregateRequest) (GetChampionAggregateResponse, error) {
	panic("not implemented")
}
