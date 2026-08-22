package profile

import "context"

type SearchProfileRequest struct {
	Query string
}

type SearchProfileResponse struct {
	Region        string
	Puuid         string
	Name          string
	Tag           string
	SummonerLevel int
	ProfileIconID int
	Rank          *Rank
}

func (s *Service) SearchProfile(ctx context.Context, req *SearchProfileRequest) (*SearchProfileResponse, error) {
	panic("bruh")
}
