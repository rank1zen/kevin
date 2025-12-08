package service

import (
	"context"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

type SearchService Service

type SearchProfileRequest struct {
	Region *riot.Region `json:"region"`
	Query  string       `json:"query"`
}

func (s *SearchService) SearchProfile(ctx context.Context, req SearchProfileRequest) ([]internal.SearchResult, error) {
	storeResults, err := s.store.Profile.SearchSummoner(ctx, req.Query)
	if err != nil {
		return nil, err
	}

	return storeResults, nil
}
