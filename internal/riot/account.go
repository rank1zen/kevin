package riot

import (
	"context"
	"fmt"
)

type AccountService service

type PUUID string

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
func (m *AccountService) GetAccountByPUUID(ctx context.Context, continent Continent, puuid string) (*Account, error) {
	endpoint := fmt.Sprintf("/riot/account/v1/accounts/by-puuid/%s", puuid)
	var account Account
	if err := m.client.makeAndDispatchRequestOnContinent(ctx, continent, endpoint, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// GetAccountByRiotID returns a account by riot id.
//
// Riot API docs: https://developer.riotgames.com/apis#account-v1/GET_getByRiotId
//
// GET /riot/account/v1/accounts/by-riot-id/{gameName}/{tagLine}
func (m *AccountService) GetAccountByRiotID(ctx context.Context, continent Continent, gameName, tagLine string) (*Account, error) {
	endpoint := fmt.Sprintf("/riot/account/v1/accounts/by-riot-id/%s/%s", gameName, tagLine)
	var account Account
	if err := m.client.makeAndDispatchRequestOnContinent(ctx, continent, endpoint, &account); err != nil {
		return nil, err
	}

	return &account, nil
}
