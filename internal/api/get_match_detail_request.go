package api

import "github.com/rank1zen/kevin/internal/riot"

type MatchDetailRequest struct {
	Region  riot.Region `json:"region"`
	MatchID string
}

func (r MatchDetailRequest) Validate() (problems map[string]string) {
	return nil
}
