package api

import (
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type GetSummonerChampionsRequest struct {
	Region riot.Region `json:"region"`
	PUUID  riot.PUUID  `json:"puuid"`
	Week   time.Time   `json:"week"`
}

func (r GetSummonerChampionsRequest) Validate() (problems map[string]string) {
	return nil
}
