package frontend

import (
	"context"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

type ProfileService Handler

type GetSummonerPageRequest struct {
	Region *riot.Region
	Name   string
	Tag    string
}

func (s *ProfileService) GetSummonerPage(ctx context.Context, req GetSummonerPageRequest) (*internal.ProfileDetail, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
	}

	storeProfile, err := s.Datasource.GetProfileDetailByRiotID(ctx, *req.Region, req.Name, req.Tag)
	if err != nil {
		return nil, err
	}

	return &storeProfile, nil
}

type UpdateProfileRequest struct {
	Region *riot.Region `json:"region"`
	Name   string       `json:"name"`
	Tag    string       `json:"tag"`
}

func (s *ProfileService) UpdateProfile(ctx context.Context, req UpdateProfileRequest) error {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
	}

	if err := s.Datasource.UpdateProfileByRiotID(ctx, *req.Region, req.Name, req.Tag); err != nil {
		return err
	}

	return nil
}

type GetMatchHistoryRequest struct {
	Region  *riot.Region `json:"region"`
	PUUID   riot.PUUID   `json:"puuid"`
	StartTS *time.Time   `json:"startTs"`
	EndTS   *time.Time   `json:"endTs"`
}

func (s *ProfileService) GetMatchHistory(ctx context.Context, req GetMatchHistoryRequest) ([]internal.SummonerMatch, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
	}

	currTime := time.Now().In(time.UTC)

	if req.StartTS == nil {
		req.StartTS = new(time.Time)
		*req.StartTS = currTime.AddDate(0, 0, -1)
	}

	if req.EndTS == nil {
		req.EndTS = new(time.Time)
		*req.EndTS = currTime
	}

	storeMatches, err := s.Datasource.GetMatchHistory(ctx, *req.Region, req.PUUID, *req.StartTS, *req.EndTS)
	if err != nil {
		return nil, err
	}

	return storeMatches, nil
}

type GetMatchDetailRequest struct {
	Region  *riot.Region `json:"region"`
	MatchID string       `json:"matchId"`
}

func (s *ProfileService) GetMatchDetail(ctx context.Context, req GetMatchDetailRequest) (*internal.MatchDetail, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
	}

	storeMatch, err := s.Datasource.GetMatchDetail(ctx, *req.Region, req.MatchID)
	if err != nil {
		return nil, err
	}

	return &storeMatch, nil
}

type GetLiveMatchRequest struct {
	Region *riot.Region `json:"region"`
	PUUID  riot.PUUID   `json:"puuid"`
}

func (s *ProfileService) GetLiveMatch(ctx context.Context, req GetLiveMatchRequest) (*internal.LiveMatch, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
	}

	match, err := s.Datasource.GetLiveMatch(ctx, *req.Region, req.PUUID)
	if err != nil {
		return nil, err
	}

	return &match, nil
}

type GetSummonerChampionsRequest struct {
	Region  *riot.Region `json:"region"`
	PUUID   riot.PUUID   `json:"puuid"`
	StartTS *time.Time   `json:"startTs"`
	EndTS   *time.Time   `json:"endTs"`
}

func (s *ProfileService) GetSummonerChampions(ctx context.Context, req GetSummonerChampionsRequest) ([]internal.SummonerChampion, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
	}

	currTime := time.Now().In(time.UTC)

	if req.StartTS == nil {
		req.StartTS = new(time.Time)
		*req.StartTS = currTime.AddDate(0, 0, -7)
	}

	if req.EndTS == nil {
		req.EndTS = new(time.Time)
		*req.EndTS = currTime
	}

	storeChamps, err := s.Datasource.GetSummonerChampions(ctx, *req.Region, req.PUUID, *req.StartTS, *req.EndTS)
	if err != nil {
		return nil, err
	}

	return storeChamps, nil
}
