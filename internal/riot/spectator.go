package riot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type CurrentGameInfo struct {
	GameID            int                      `json:"gameId"`
	GameType          string                   `json:"gameType"`
	GameStartTime     int64                    `json:"gameStartTime"`
	MapID             int64                    `json:"mapId"`
	GameLength        int64                    `json:"gameLength"`
	PlatformID        string                   `json:"platformId"`
	GameMode          string                   `json:"gameMode"`
	BannedChampions   []BannedChampion         `json:"bannedChampions"`
	GameQueueConfigID int64                    `json:"gameQueueConfigId"`
	Observers         Observer                 `json:"observers"`
	Participants      []CurrentGameParticipant `json:"participants"`
}

type BannedChampion struct {
	PickTurn   int `json:"pickTurn"`
	ChampionID int `json:"championId"`
	TeamID     int `json:"teamId"`
}

type Observer struct {
	EncryptionKey string `json:"encryptionKey"`
}

type CurrentGameParticipant struct {
	ChampionID               int                       `json:"championId"`
	Perks                    SpectatorPerks            `json:"perks"`
	ProfileIconID            int                       `json:"profileIconId"`
	Bot                      bool                      `json:"bot"`
	TeamID                   int                       `json:"teamId"`
	SummonerID               string                    `json:"summonerId"`
	PUUID                    string                    `json:"puuid"`
	Spell1ID                 int                       `json:"spell1Id"`
	Spell2ID                 int                       `json:"spell2Id"`
	GameCustomizationObjects []GameCustomizationObject `json:"gameCustomizationObjects"`
}

type SpectatorPerks struct {
	PerkIDs      []int `json:"perkIds"`
	PerkStyle    int   `json:"perkStyle"`
	PerkSubStyle int   `json:"perkSubStyle"`
}

type GameCustomizationObject struct {
	Category string `json:"category"`
	Content  string `json:"content"`
}

// GetCurrentGameInfoByPUUID returns current game information for the given puuid.
//
// Riot API docs: https://developer.riotgames.com/apis#spectator-v5/GET_getCurrentGameInfoByPuuid
//
// GET /lol/spectator/v5/active-games/by-summoner/{encryptedPUUID}
func (c *Client) GetCurrentGameInfoByPUUID(ctx context.Context, platform, puuid string) (*CurrentGameInfo, error) {
	u := platformHost(platform)
	path := fmt.Sprintf("/lol/spectator/v5/active-games/by-summoner/%s", puuid)
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
	var m CurrentGameInfo
	err = json.NewDecoder(resp.Body).Decode(&m)
	return &m, err
}
