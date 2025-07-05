package riot

import (
	"context"
	"fmt"
)

type SpectatorService service

type LiveMatch struct {
	GameID            int                    `json:"gameId"`
	GameType          string                 `json:"gameType"`
	GameStartTime     int64                  `json:"gameStartTime"`
	MapID             int64                  `json:"mapId"`
	GameLength        int64                  `json:"gameLength"`
	PlatformID        string                 `json:"platformId"`
	GameMode          string                 `json:"gameMode"`
	BannedChampions   []LiveBannedChampion   `json:"bannedChampions"`
	GameQueueConfigID int64                  `json:"gameQueueConfigId"`
	Observers         LiveObserver           `json:"observers"`
	Participants      []LiveMatchParticipant `json:"participants"`
}

type LiveBannedChampion struct {
	PickTurn   int `json:"pickTurn"`
	ChampionID int `json:"championId"`
	TeamID     int `json:"teamId"`
}

type LiveObserver struct {
	EncryptionKey string `json:"encryptionKey"`
}

type LiveMatchParticipant struct {
	ChampionID               int                           `json:"championId"`
	Perks                    LivePerks                     `json:"perks"`
	ProfileIconID            int                           `json:"profileIconId"`
	Bot                      bool                          `json:"bot"`
	TeamID                   int                           `json:"teamId"`
	SummonerID               string                        `json:"summonerId"`
	PUUID                    string                        `json:"puuid"`
	Spell1ID                 int                           `json:"spell1Id"`
	Spell2ID                 int                           `json:"spell2Id"`
	GameCustomizationObjects []LiveGameCustomizationObject `json:"gameCustomizationObjects"`
}

type LivePerks struct {
	PerkIDs      []int `json:"perkIds"`
	PerkStyle    int   `json:"perkStyle"`
	PerkSubStyle int   `json:"perkSubStyle"`
}

type LiveGameCustomizationObject struct {
	Category string `json:"category"`
	Content  string `json:"content"`
}

// GetLiveMatch returns current game information for the given puuid.
//
// Riot API docs: https://developer.riotgames.com/apis#spectator-v5/GET_getCurrentGameInfoByPuuid
//
// GET /lol/spectator/v5/active-games/by-summoner/{encryptedPUUID}
func (m *SpectatorService) GetLiveMatch(ctx context.Context, region Region, puuid string) (*LiveMatch, error) {
	path := fmt.Sprintf("/lol/spectator/v5/active-games/by-summoner/%s", puuid)

	var match LiveMatch
	if err := m.client.makeAndDispatchRequest(ctx, region, path, &match); err != nil {
		return nil, err
	}

	return &match, nil
}
