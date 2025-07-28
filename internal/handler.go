package internal

import (
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type MatchHistoryRequest struct {
	Region riot.Region `json:"region"`

	PUUID riot.PUUID `json:"puuid"`

	// Date should be the start of the day. The request will fetch all
	// matches played on the day.
	Date time.Time `json:"date"`
}

func (r MatchHistoryRequest) Validate() (problems map[string]string) {
	problems = make(map[string]string)

	validatePUUID(problems, r.PUUID)

	if r.Date.Hour() != 0 || r.Date.Minute() != 0 || r.Date.Second() != 0 {
		problems["date"] = "date needs to be start of day"
	}

	return problems
}

type SummonerChampionRequest struct {
	Region riot.Region `json:"region"`
	PUUID  riot.PUUID  `json:"puuid"`
	Week   time.Time   `json:"week"`
}

func (r SummonerChampionRequest) Validate() (problems map[string]string) {
	problems = make(map[string]string)

	validatePUUID(problems, r.PUUID)

	if r.Week.Hour() != 0 || r.Week.Minute() != 0 || r.Week.Second() != 0 {
		problems["date"] = "date needs to be start of day"
	}

	return problems
}

type SummonerPage struct {
	// Region specifies a riot region. All results in the page belong to
	// this region.
	Region riot.Region

	PUUID riot.PUUID

	Name, Tag string

	LastUpdated time.Time

	Rank *Rank

	MatchHistoryRequest []MatchHistoryRequest
}
