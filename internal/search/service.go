package search

import (
	"context"
)

type Service struct {
	store Store
}

func (s *Service) SearchSummoner(ctx context.Context, q string) ([]SearchResult, error) {
	return s.store.SearchSummoner(ctx, q)
}

func (s *Service) SearchByNameTag(ctx context.Context, name, tag string) ([]Profile, error) {
	return s.store.SearchByNameTag(ctx, name, tag)
}
