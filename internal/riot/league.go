package riot

import (
	"context"
	"fmt"
)

type LeagueService service

const (
	QueueTypeRankedSolo5x5 = "RANKED_SOLO_5x5"
	QueueTypeRankedFlexSR = "RANKED_FLEX_SR"
	QueueTypeRankedFlexTT = "RANKED_FLEX_TT"
)

type Tier string

func (t Tier) String() string {
	var tier string
	switch t {
	case TierIron:
		tier = "Iron"
	case TierBronze:
		tier = "Bronze"
	case TierSilver:
		tier = "Silver"
	case TierGold:
		tier = "Gold"
	case TierPlatinum:
		tier = "Platinum"
	case TierEmerald:
		tier = "Emerald"
	case TierDiamond:
		tier = "Diamond"
	case TierMaster:
		tier = "Master"
	case TierGrandmaster:
		tier = "Grandmaster"
	case TierChallenger:
		tier = "Challenger"
	}
	return tier
}

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

type Division string

func (r Division) String() string {
	var tier string
	switch r {
	case Division1:
		tier = "I"
	case Division2:
		tier = "II"
	case Division3:
		tier = "III"
	case Division4:
		tier = "IV"
	}
	return tier
}

const (
	Division1 Division = "I"
	Division2 Division = "II"
	Division3 Division = "III"
	Division4 Division = "IV"
)

type LeagueList []LeagueEntry

type LeagueEntry struct {
	Division         Division        `json:"rank"`
	FreshBlood   bool        `json:"freshBlood"`
	HotStreak    bool        `json:"hotStreak"`
	Inactive     bool        `json:"inactive"`
	LeagueID     string      `json:"leagueId"`
	LeaguePoints int         `json:"leaguePoints"`
	Losses       int         `json:"losses"`
	MiniSeries   *LeagueMiniSeries `json:"miniSeries"`
	QueueType    string      `json:"queueType"`
	SummonerID   string      `json:"summonerId"`
	Tier         Tier        `json:"tier"`
	Veteran      bool        `json:"veteran"`
	Wins         int         `json:"wins"`
}

type LeagueMiniSeries struct {
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
func (m *LeagueService) GetLeagueEntriesForSummoner(ctx context.Context, region Region, summonerID string) (LeagueList, error) {
	path := fmt.Sprintf("/lol/league/v4/entries/by-summoner/%v", summonerID)

	var entries LeagueList
	if err := m.client.makeAndDispatchRequest(ctx, region, path, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// GetLeagueEntries returns league entries in all queues for a given summoner ID.
//
// Riot API docs: https://developer.riotgames.com/apis#league-v4/GET_getLeagueEntriesByPUUID
//
// GET /lol/league/v4/entries/by-puuid/{encryptedPUUID}
func (m *LeagueService) GetLeagueEntriesByPUUID(ctx context.Context, region Region, puuid string) (LeagueList, error) {
	path := fmt.Sprintf("/lol/league/v4/entries/by-puuid/%v", puuid)

	var entries LeagueList
	if err := m.client.makeAndDispatchRequest(ctx, region, path, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}
