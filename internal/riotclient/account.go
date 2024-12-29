package riotclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Account struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

// Get account by riot id
func (c *Client) AccountGetByRiotId(ctx context.Context, gameName, tagLine string) (*Account, error) {
	u := fmt.Sprintf(defaultAmerBaseURL+"/riot/account/v1/accounts/by-riot-id/%s/%s", gameName, tagLine)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	body, err := execute(ctx, req)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var m Account
	err = json.NewDecoder(body).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("riot: json error (%v)", err)
	}

	return &m, nil
}

// Get account by puuid
func (c *Client) AccountGetByPuuid(ctx context.Context, puuid string) (*Account, error) {
	u := fmt.Sprintf(defaultAmerBaseURL+"/riot/account/v1/accounts/by-puuid/%s", puuid)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	body, err := execute(ctx, req)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var m Account
	err = json.NewDecoder(body).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("riot: json error (%v)", err)
	}

	return &m, nil
}
