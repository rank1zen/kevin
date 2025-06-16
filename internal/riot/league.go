package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const QueueTypeRankedSolo5x5 = "RANKED_SOLO_5x5"

type Tier string

const (
	TierIron        Tier = "IRON"
	TierBronze      Tier = "BRONZE"
	TierSilver      Tier = "SILVER"
	TierGold        Tier = "GOLD"
	TierPlatinum    Tier = "PLATINUM"
	TierEmerald     Tier = "EMERALD"
	TierDiamond     Tier = "DIAMOND"
	TierMaster      Tier = "MASTER"
	TierGrandmaster Tier = "GRANDMASTER"
	TierChallenger  Tier = "CHALLENGER"
)

type Rank string

const (
	Rank1 Rank = "I"
	Rank2 Rank = "II"
	Rank3 Rank = "III"
	Rank4 Rank = "IV"
)

type LeagueEntry struct {
	FreshBlood   bool        `json:"freshBlood"`
	HotStreak    bool        `json:"hotStreak"`
	Inactive     bool        `json:"inactive"`
	LeagueID     string      `json:"leagueId"`
	LeaguePoints int         `json:"leaguePoints"`
	Losses       int         `json:"losses"`
	MiniSeries   *MiniSeries `json:"miniSeries"`
	QueueType    string      `json:"queueType"`
	Rank         Rank        `json:"rank"`
	SummonerID   string      `json:"summonerId"`
	Tier         Tier        `json:"tier"`
	Veteran      bool        `json:"veteran"`
	Wins         int         `json:"wins"`
}

type MiniSeries struct {
	Losses   int    `json:"losses"`
	Progress string `json:"progess"`
	Target   int    `json:"target"`
	Wins     int    `json:"wins"`
}

// GetLeagueEntriesForSummoner returns league entries in all queues for a given summoner ID.
//
// Riot API docs: https://developer.riotgames.com/apis#league-v4/GET_getLeagueEntriesForSummoner
//
// GET /lol/league/v4/entries/by-summoner/{encryptedSummonerId}
func (c *Client) GetLeagueEntriesForSummoner(ctx context.Context, platform, summonerID string) (entries []*LeagueEntry, err error) {
	u := platformHost(platform)
	path := fmt.Sprintf("/lol/league/v4/entries/by-summoner/%v", summonerID)
	u = u.JoinPath(path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Riot-Token", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&entries)
	return entries, err
}

// GetLeagueEntries returns league entries in all queues for a given summoner ID.
//
// Riot API docs: https://developer.riotgames.com/apis#league-v4/GET_getLeagueEntriesByPUUID
func (c *Client) GetLeagueEntries(ctx context.Context, platform, puuid string) (entries []*LeagueEntry, err error) {
	u := platformHost(platform)
	path := fmt.Sprintf("/lol/league/v4/entries/by-puuid/%v", puuid)
	u = u.JoinPath(path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Riot-Token", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&entries)
	return entries, err
}
