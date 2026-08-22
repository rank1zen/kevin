package profile

import (
	"context"
)

type UpdateProfileRequest struct {
	Region string
	Name   string
	Tag    string
}

type UpdateProfileResponse struct {
	Region        string
	Puuid         string
	Name          string
	Tag           string
	SummonerLevel int
	ProfileIconID int
	Rank          *Rank
}

func (s *Service) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) (*UpdateProfileResponse, error) {
	panic("bruh")
}
