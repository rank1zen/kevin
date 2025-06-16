package frontend

import (
	"encoding/json"

	"github.com/rank1zen/kevin/internal/riot"
)

type EndpointProvider struct {}

func (ep EndpointProvider) GetMatchHistory(region riot.Region, puuid riot.PUUID, index int) (string, []byte) {
	path := "/summoner/matchlist"

	req := MatchHistoryRequest{
		Region: region,
		PUUID:  puuid,
		Date:   GetDay(index),
	}

	bytes, _ := json.Marshal(req)

	return path, bytes
}

func (ep EndpointProvider) GetLiveMatch(region riot.Region, puuid riot.PUUID) (string, []byte) {
	path := "/summoner/live"

	req := LiveMatchRequest{
		Region: region,
		PUUID:  puuid,
	}

	bytes, _ := json.Marshal(req)

	return path, bytes
}

func (ep EndpointProvider) GetChampionList(region riot.Region, puuid riot.PUUID) (string, []byte) {
	path := "/summoner/champions"

	req := ZGetSummonerChampionsRequest{
		Region: region,
		PUUID:  puuid,
		Week:   GetCurrentWeek(),
	}

	bytes, _ := json.Marshal(req)

	return path, bytes
}
