package riotclient

import (
	"context"
	"encoding/json"
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

// Get a summoner by PUUID.
func (c *Client) GetSummonerByPuuid(ctx context.Context, puuid string) (*Summoner, error) {
	u := fmt.Sprintf(defaultNaBaseURL+"/lol/summoner/v4/summoners/by-puuid/%v", puuid)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	body, err := execute(ctx, req)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var m *Summoner
	err = json.NewDecoder(body).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("riot: json error (%v)", err)
	}

	return m, nil
}
