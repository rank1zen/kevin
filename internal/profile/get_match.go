package profile

import "context"

type GetMatchRequest struct {
	Region  string
	MatchID string
}

type GetMatchResponse struct {
}

func (s *Service) GetMatch(ctx context.Context, req GetMatchRequest) (GetMatchResponse, error) {
	panic("not implemented")
}
