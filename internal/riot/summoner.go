package riot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Summoner struct {
	AccountId     string `json:"accountId"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	Id            string `json:"id"`
	Puuid         string `json:"puuid"`
	SummonerLevel int64  `json:"summonerLevel"`
}

// GetSummonerByPuuid returns a summoner by PUUID.
//
// Riot API docs: https://developer.riotgames.com/apis#summoner-v4/GET_getByPUUID
//
// GET /lol/summoner/v4/summoners/by-puuid/{encryptedPUUID}
func (c *Client) GetSummonerByPuuid(ctx context.Context, platform, puuid string) (*Summoner, error) {
	u := platformHost(platform)
	path := fmt.Sprintf("/lol/summoner/v4/summoners/by-puuid/%s", puuid)
	u = u.JoinPath(path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Riot-Token", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	var summoner Summoner
	err = json.NewDecoder(resp.Body).Decode(&summoner)
	return &summoner, err
}
