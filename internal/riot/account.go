package riot

import (
	"context"
	"fmt"

	"github.com/rank1zen/kevin/internal/riot/internal"
)

// AccountService is the ACCOUNT-V1 API.
//
// Riot API docs: https://developer.riotgames.com/apis#account-v1
type AccountService service

// PUUID is a 78 character global identifier for a Riot account.
type PUUID string

func (id PUUID) String() string {
	return string(id)
}

type Account struct {
	PUUID    PUUID  `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

// GetAccountByPuuid returns a account by puuid.
//
// Riot API docs: https://developer.riotgames.com/apis#account-v1/GET_getByPuuid
//
// GET /riot/account/v1/accounts/by-puuid/{puuid}
func (m *AccountService) GetAccountByPUUID(ctx context.Context, region Region, puuid string) (*Account, error) {
	endpoint := fmt.Sprintf("/riot/account/v1/accounts/by-puuid/%s", puuid)

	req := &internal.Request{
		BaseURL:  region.host(),
		Endpoint: endpoint,
		APIKey:   m.client.apiKey,
	}

	if m.client.baseURL != "" {
		req.BaseURL = m.client.baseURL
	}

	var account Account
	if err := m.client.internals.DispatchRequest(ctx, req, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// GetAccountByRiotID returns a account by riot id.
//
// Riot API docs: https://developer.riotgames.com/apis#account-v1/GET_getByRiotId
//
// GET /riot/account/v1/accounts/by-riot-id/{gameName}/{tagLine}
func (m *AccountService) GetAccountByRiotID(ctx context.Context, region Region, gameName, tagLine string) (*Account, error) {
	endpoint := fmt.Sprintf("/riot/account/v1/accounts/by-riot-id/%s/%s", gameName, tagLine)

	req := &internal.Request{
		BaseURL:  region.host(),
		Endpoint: endpoint,
		APIKey:   m.client.apiKey,
	}

	if m.client.baseURL != "" {
		req.BaseURL = m.client.baseURL
	}

	var account Account
	if err := m.client.internals.DispatchRequest(ctx, req, &account); err != nil {
		return nil, err
	}

	return &account, nil
}
