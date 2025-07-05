package riot

import (
	"context"
	"fmt"
)

type SummonerService service

type Summoner struct {
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	PUUID         string `json:"puuid"`
	SummonerLevel int64  `json:"summonerLevel"`
}

// GetSummonerByPuuid returns a summoner by PUUID.
//
// Riot API docs: https://developer.riotgames.com/apis#summoner-v4/GET_getByPUUID
//
// GET /lol/summoner/v4/summoners/by-puuid/{encryptedPUUID}
func (m *SummonerService) GetSummoner(ctx context.Context, region Region, puuid string) (*Summoner, error) {
	path := fmt.Sprintf("/lol/summoner/v4/summoners/by-puuid/%s", puuid)

	var summoner Summoner
	if err := m.client.makeAndDispatchRequest(ctx, region, path, &summoner); err != nil {
		return nil, err
	}

	return &summoner, nil
}
