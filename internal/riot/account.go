package riot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Account struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

// GetAccountByPuuid returns a account by puuid.
//
// Riot API docs: https://developer.riotgames.com/apis#account-v1/GET_getByPuuid
//
// GET /riot/account/v1/accounts/by-puuid/{puuid}
func (c *Client) GetAccountByPuuid(ctx context.Context, region, puuid string) (*Account, error) {
	u := regionHost(region)
	path := fmt.Sprintf("/riot/account/v1/accounts/by-puuid/%s", puuid)
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
	var account Account
	err = json.NewDecoder(resp.Body).Decode(&account)
	return &account, err
}

// GetAccountByRiotID returns a account by riot id.
//
// Riot API docs: https://developer.riotgames.com/apis#account-v1/GET_getByRiotId
//
// GET /riot/account/v1/accounts/by-riot-id/{gameName}/{tagLine}
func (c *Client) GetAccountByRiotID(ctx context.Context, region, gameName, tagLine string) (*Account, error) {
	u := regionHost(region)
	path := fmt.Sprintf("/riot/account/v1/accounts/by-riot-id/%s/%s", gameName, tagLine)
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
	var account Account
	err = json.NewDecoder(resp.Body).Decode(&account)
	return &account, err
}
