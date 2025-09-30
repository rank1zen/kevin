package api

import "github.com/rank1zen/kevin/internal/riot"

type UpdateSummonerRequest struct {
	Region riot.Region `json:"region"`
	Name   string      `json:"name"`
	Tag    string      `json:"tag"`
}

func (r UpdateSummonerRequest) Validate() (problems map[string]string) {
	return nil
}
