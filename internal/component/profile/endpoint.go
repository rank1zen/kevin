package profile

import "github.com/rank1zen/kevin/internal/riot"

type EndpointProvider interface {
	GetMatchHistory(region riot.Region, puuid riot.PUUID, index int) (path string, json []byte)

	GetLiveMatch(region riot.Region, puuid riot.PUUID) (path string, json []byte)

	GetChampionList(region riot.Region, puuid riot.PUUID) (path string, json []byte)
}
