package profile

import (
	"context"
	"time"
)

type GetMatchHistoryRequest struct {
	PUUID               string
	TimestampRangeStart time.Time
	TimestampRangeEnd   time.Time
	PageSize            int
	PageToken           string
}

type GetMatchHistoryResponse struct {
	Matches       []Match
	NextPageToken string
}

func (s *Service) GetMatchHistory(ctx context.Context, req GetMatchHistoryRequest) (GetMatchHistoryResponse, error) {
	panic("not implemented")
}
