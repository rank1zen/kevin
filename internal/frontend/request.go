package frontend

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
	problems = make(map[string]string)

	validatePUUID(problems, r.PUUID)

	return problems
}

type MatchDetailRequest struct {
	Region  riot.Region `json:"region"`
	MatchID string
}

func (r MatchDetailRequest) Validate() (problems map[string]string) {
	return nil
}

type GetSummonerChampionsRequest struct {
	Region riot.Region `json:"region"`
	PUUID  riot.PUUID  `json:"puuid"`
	Week   time.Time   `json:"week"`
}

func (r GetSummonerChampionsRequest) Validate() (problems map[string]string) {
	problems = make(map[string]string)

	validatePUUID(problems, r.PUUID)

	if r.Week.Hour() != 0 || r.Week.Minute() != 0 || r.Week.Second() != 0 {
		problems["date"] = "date needs to be start of day"
	}

	return problems
}

type LiveMatchRequest struct {
	Region riot.Region `json:"region"`
	PUUID  riot.PUUID  `json:"puuid"`
}

func (r LiveMatchRequest) Validate() (problems map[string]string) {
	problems = make(map[string]string)
	validatePUUID(problems, r.PUUID)
	return problems
}

type UpdateSummonerRequest struct {
	Region riot.Region `json:"region"`
	Name   string      `json:"name"`
	Tag    string      `json:"tag"`
}

func (r UpdateSummonerRequest) Validate() (problems map[string]string) {
	return nil
}
