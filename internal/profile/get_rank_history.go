package profile

import (
	"context"
	"time"
)

type GetRankHistoryRequest struct {
	PUUID     string
	StartDate time.Time
	EndDate   time.Time
}

type GetRankHistoryResponse struct {
}

func (s *Service) GetRankHistory(ctx context.Context, req GetRankHistoryRequest) (GetRankHistoryResponse, error) {
	panic("not implemented")
}
