package api

import "github.com/rank1zen/kevin/internal/riot"

type LiveMatchRequest struct {
	Region riot.Region `json:"region"`
	PUUID  riot.PUUID  `json:"puuid"`
}

func (r LiveMatchRequest) Validate() (problems map[string]string) {
	return nil
}
