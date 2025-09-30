package api

import (
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type MatchHistoryRequest struct {
	Region  riot.Region `json:"region"`
	PUUID   riot.PUUID  `json:"puuid"`
	StartTS time.Time   `json:"startTs"`
	EndTS   time.Time   `json:"endTs"`
}

func (r MatchHistoryRequest) Validate() (problems map[string]string) {
	return nil
}
