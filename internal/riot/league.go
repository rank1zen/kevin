package riot

import (
	"context"
	"fmt"

	"github.com/rank1zen/kevin/internal/riot/internal"
)

type LeagueService service

const (
	QueueTypeRankedSolo5x5 = "RANKED_SOLO_5x5"
	QueueTypeRankedFlexSR  = "RANKED_FLEX_SR"
	QueueTypeRankedFlexTT  = "RANKED_FLEX_TT"
)

type LeagueList []LeagueEntry

type LeagueEntry struct {
	Division     string            `json:"rank"`
	FreshBlood   bool              `json:"freshBlood"`
	HotStreak    bool              `json:"hotStreak"`
	Inactive     bool              `json:"inactive"`
	LeagueID     string            `json:"leagueId"`
	LeaguePoints int               `json:"leaguePoints"`
	Losses       int               `json:"losses"`
	MiniSeries   *LeagueMiniSeries `json:"miniSeries"`
	QueueType    string            `json:"queueType"`
	SummonerID   string            `json:"summonerId"`
	Tier         string            `json:"tier"`
	Veteran      bool              `json:"veteran"`
	Wins         int               `json:"wins"`
}

type LeagueMiniSeries struct {
	Losses   int    `json:"losses"`
	Progress string `json:"progress"`
	Target   int    `json:"target"`
	Wins     int    `json:"wins"`
}

// GetLeagueEntries returns league entries in all queues for a given summoner ID.
//
// Riot API docs: https://developer.riotgames.com/apis#league-v4/GET_getLeagueEntriesByPUUID
//
// GET /lol/league/v4/entries/by-puuid/{encryptedPUUID}
func (m *LeagueService) GetLeagueEntriesByPUUID(ctx context.Context, region string, puuid string) (LeagueList, error) {
	endpoint := fmt.Sprintf("/lol/league/v4/entries/by-puuid/%v", puuid)

	req := &internal.Request{
		BaseURL:  regionHost(region),
		Endpoint: endpoint,
		APIKey:   m.client.apiKey,
	}

	if m.client.baseURL != "" {
		req.BaseURL = m.client.baseURL
	}

	var entries LeagueList
	if err := m.client.internals.DispatchRequest(ctx, req, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// GetChallengerLeague returns the challenger league for a given queue.
//
// Riot API docs: https://developer.riotgames.com/apis#league-v4/GET_getChallengerLeague
func (m *LeagueService) GetChallengerLeague(ctx context.Context, region string, queue string) (*LeagueList2, error) {
	endpoint := fmt.Sprintf("/lol/league/v4/challengerleagues/by-queue/%s", queue)

	req := &internal.Request{
		BaseURL:  regionHost(region),
		Endpoint: endpoint,
		APIKey:   m.client.apiKey,
	}

	if m.client.baseURL != "" {
		req.BaseURL = m.client.baseURL
	}

	var entries LeagueList2
	if err := m.client.internals.DispatchRequest(ctx, req, &entries); err != nil {
		return nil, err
	}

	return &entries, nil
}

type LeagueList2 struct {
	LeagueID string       `json:"leagueId"`
	Entries  []LeagueItem `json:"entries"`
	Tier     string       `json:"tier"`
	Name     string       `json:"name"`
	Queue    string       `json:"queue"`
}

type LeagueItem struct {
	FreshBlood   bool              `json:"freshBlood"`
	HotStreak    bool              `json:"hotStreak"`
	Inactive     bool              `json:"inactive"`
	Rank         string            `json:"rank"`
	LeaguePoints int               `json:"leaguePoints"`
	Losses       int               `json:"losses"`
	MiniSeries   *LeagueMiniSeries `json:"miniSeries"`
	Veteran      bool              `json:"veteran"`
	Wins         int               `json:"wins"`
	PUUID        string            `json:"puuid"`
}
