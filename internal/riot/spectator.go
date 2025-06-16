package riot

import (
	"context"
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

type SpectatorPerks struct {
	PerkIds      []int `json:"perkIds"`
	PerkStyle    int   `json:"perkStyle"`
	PerkSubStyle int   `json:"perkSubStyle"`
}

type SpectatorCustomizationObject struct {
	Category string `json:"category"`
	Content  string `json:"content"`
}

// GetCurrentGameInfoByPuuid returns current game information for the given puuid.
//
// Riot API docs: https://developer.riotgames.com/apis#spectator-v5/GET_getCurrentGameInfoByPuuid
//
// GET /lol/spectator/v5/active-games/by-summoner/{encryptedPUUID}
func (c *Client) GetCurrentGameInfoByPuuid(ctx context.Context, platform, puuid string) (m SpectatorCurrentGameInfo, err error) {
	panic("not implemented")
}
