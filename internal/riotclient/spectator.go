package riotclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SpectatorCurrentGameInfo struct {
	GameId            int                               `json:"gameId"`
	GameType          string                            `json:"gameType"`
	GameStartTime     int64                             `json:"gameStartTime"`
	MapId             int64                             `json:"mapId"`
	GameLength        int64                             `json:"gameLength"`
	PlatformId        string                            `json:"platformId"`
	GameMode          string                            `json:"gameMode"`
	BannedChampions   []SpectatorBannedChampion         `json:"bannedChampions"`
	GameQueueConfigId int64                             `json:"gameQueueConfigId"`
	Observers         SpectatorObserver                 `json:"observers"`
	Participants      []SpectatorCurrentGameParticipant `json:"participants"`
}

type SpectatorBannedChampion struct {
	PickTurn   int `json:"pickTurn"`
	ChampionId int `json:"championId"`
	TeamId     int `json:"teamId"`
}

type SpectatorObserver struct {
	EncryptionKey string `json:"encryptionKey"`
}

type SpectatorCurrentGameParticipant struct {
	ChampionId               int                            `json:"encryptionKey"`
	Perks                    SpectatorPerks                 `json:"perks"`
	ProfileIconId            int                            `json:"profileIconId"`
	Bot                      bool                           `json:"bot"`
	TeamId                   int                            `json:"teamId"`
	SummonerId               string                         `json:"summonerId"`
	Puuid                    string                         `json:"puuid"`
	Spell1Id                 int                            `json:"spell1Id"`
	Spell2Id                 int                            `json:"spell2Id"`
	GameCustomizationObjects []SpectatorCustomizationObject `json:"gameCustomizationObjects"`
}

const PerkKeystone = 0
const PerkSecondary = 1

type SpectatorPerks struct {
	PerkIds      []int `json:"perkIds"`
	PerkStyle    int   `json:"perkStyle"`
	PerkSubStyle int   `json:"perkSubStyle"`
}

type SpectatorCustomizationObject struct {
	Category string `json:"category"`
	Content  string `json:"content"`
}

func (c *Client) GetCurrentGameInfoByPuuid(ctx context.Context, puuid string) (SpectatorCurrentGameInfo, error) {
	var m SpectatorCurrentGameInfo

	u := fmt.Sprintf(defaultNaBaseURL+"/lol/spectator/v5/active-games/by-summoner/%v", puuid)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return m, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	body, err := execute(ctx, req)
	if err != nil {
		return m, err
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&m)
	if err != nil {
		return m, err
	}

	return m, nil
}
