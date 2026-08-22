package profile

import (
	"context"
)

type GetProfileRequest struct {
	Region string
	Name   string
	Tag    string
}

type GetProfileResponse struct {
	Region        string
	Puuid         string
	Name          string
	Tag           string
	SummonerLevel int
	ProfileIconID int
	Rank          *Rank
}

func (s *Service) GetProfile(ctx context.Context, req GetProfileRequest) (GetProfileResponse, error) {
	panic("not implemented")
}
