package riotclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const QueueTypeRankedSolo5x5 = "RANKED_SOLO_5x5"

type LeagueEntry struct {
	LeagueId string `json:"leagueId"`

	SummonerId string `json:"summonerId"`

	QueueType string `json:"queueType"`

	// Tier is the summoner's tier, e.g., 'Diamond'
	Tier string `json:"tier"`

	// Rank is the summoner's division within a tier, e.g., 'IV'
	Rank string `json:"rank"`

	LeaguePoints int `json:"leaguePoints"`

	Wins int `json:"wins"`

	Losses int `json:"losses"`

	HotStreak bool `json:"hotStreak"`

	Veteran bool `json:"veteran"`

	FreshBlood bool `json:"freshBlood"`

	Inactive bool `json:"inactive"`

	MiniSeries *MiniSeries `json:"miniSeries"`
}

type LeagueEntryList []*LeagueEntry

type MiniSeries struct {
	Losses   int    `json:"losses"`
	Progress string `json:"progess"`
	Target   int    `json:"target"`
	Wins     int    `json:"wins"`
}

// Get league entries in all queues for a given summoner ID.
func (c *Client) GetLeagueEntriesForSummoner(ctx context.Context, summonerID string) (LeagueEntryList, error) {
	u := fmt.Sprintf(defaultNaBaseURL+"/lol/league/v4/entries/by-summoner/%v", summonerID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	body, err := execute(ctx, req)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var m LeagueEntryList
	err = json.NewDecoder(body).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("riot: json error (%v)", err)
	}

	return m, nil
}
