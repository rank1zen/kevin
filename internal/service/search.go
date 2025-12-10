package service

import (
	"context"
	"strings"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

type SearchService Service

type SearchProfileRequest struct {
	Region *riot.Region `json:"region"`
	Query  string       `json:"query"`
}

func (s *SearchService) SearchProfile(ctx context.Context, req SearchProfileRequest) (*SearchResult, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
	}

	var name, tag string

	if index := strings.Index(req.Query, "#"); index != -1 {
		name = req.Query[:index]
		tag = req.Query[index+1:]
	} else {
		name = req.Query
		tag = string(*req.Region)
	}

	storeResults, err := s.store.Profile.SearchSummoner(ctx, req.Query)
	if err != nil {
		return nil, err
	}

	result := SearchResult{
		Region:   *req.Region,
		Name:     name,
		Tag:      tag,
		Profiles: []ProfileSearchResult{},
	}

	for _, storeResult := range storeResults {
		result.Profiles = append(result.Profiles, ProfileSearchResult{
			Region: *req.Region,
			Name:   name,
			Tag:    tag,
			PUUID:  storeResult.PUUID,
			Rank:   nil, // TODO: Implement rank fetching logic
		})
	}

	return &result, nil
}

type SearchResult struct {
	Region riot.Region

	// Name and Tag is the name#tag that is used for search.
	Name, Tag string

	Profiles []ProfileSearchResult
}

type ProfileSearchResult struct {
	Region riot.Region

	PUUID riot.PUUID

	Name, Tag string

	// Rank is the summoners most recent rank record in store.
	Rank *internal.RankStatus
}
