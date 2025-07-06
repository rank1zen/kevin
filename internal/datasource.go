package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

var (
	// ErrSummonerNotFound indicates the summoner associated with some
	// puuid or name#tagline does not exist.
	ErrSummonerDoesNotExist = errors.New("summoner does not exist")

	// ErrNoLiveMatch indicates a summoner is not in a game.
	ErrNoLiveMatch = errors.New("no live match")
)

// Datasource manages business operations for the frontend. Region parameters
// specify the region to search.
type Datasource struct {
	// probably want to cache something

	riot *riot.Client

	store Store
}

func NewDatasource(client *riot.Client, store Store) *Datasource {
	return &Datasource{client, store}
}

func (ds *Datasource) ZUpdateMatchHistory(ctx context.Context, region riot.Region, puuid string, date time.Time) error {
	startTime := date.Truncate(24*time.Hour)
	endTime := startTime.Add(24 * time.Hour).Add(-1*time.Second)

	query := fmt.Sprintf("queue=420&startTime=%d&endTime=%d", startTime.Unix(), endTime.Unix())

	continent := riot.RegionToContinent(region)

	ids, err := ds.riot.Match.GetMatchList(ctx, continent, puuid, query)
	if err != nil {
		return fmt.Errorf("fetching ids: %w", err)
	}

	matchIDs := []string{}
	for _, id := range ids {
		matchIDs = append(matchIDs, id)
	}

	newIDs, err := ds.store.GetNewMatchIDs(ctx, matchIDs)
	if err != nil {
		return err
	}

	for _, id := range newIDs {
		if err := ds.recordMatch(ctx, continent, id); err != nil {
			return err
		}
	}

	return nil
}

// GetLiveMatch fetches riot for the a live match. If summoner is not in
// a game, return ErrNoLiveMatch.
func (ds *Datasource) GetLiveMatch(ctx context.Context, region riot.Region, puuid string) (LiveMatch, error) {
	riotGame, err := ds.riot.Spectator.GetLiveMatch(ctx, region, puuid)
	if err != nil {
		return LiveMatch{}, err
	}

	return NewLiveMatch(WithRiotLiveMatch(riotGame)), nil
}

// GetRiotName returns the Riot ID (name#tag) associated with puuid.
func (ds *Datasource) GetRiotName(ctx context.Context, puuid string) (name, tag string, err error) {
	// Using AMER for now since puuid is globally unique ...
	account, err := ds.riot.Account.GetAccountByPUUID(ctx, riot.ContinentAmericas, puuid)
	if err != nil {
		return "", "", err
	}

	return account.GameName, account.TagLine, nil
}

func (ds *Datasource) GetPUUID(ctx context.Context, name, tagline string) (puuid string, err error) {
	// Using AMER for now...
	account, err := ds.riot.Account.GetAccountByRiotID(ctx, riot.ContinentAmericas, name, tagline)
	if err != nil {
		return "", err
	}

	return account.PUUID, nil
}

// UpdateMatchHistory syncs the matchlist of summoner with puuid from start to
// start + count.
func (ds *Datasource) UpdateMatchHistory(ctx context.Context, region riot.Region, puuid string, start, count int) error {
	query := fmt.Sprintf("queue=420&start=%d&count=%d", start, count)

	continent := riot.RegionToContinent(region)

	ids, err := ds.riot.Match.GetMatchList(ctx, continent, puuid, query)
	if err != nil {
		return fmt.Errorf("fetching ids: %w", err)
	}

	matchIDs := []string{}
	for _, id := range ids {
		matchIDs = append(matchIDs, id)
	}

	newIDs, err := ds.store.GetNewMatchIDs(ctx, matchIDs)
	if err != nil {
		return err
	}

	for _, id := range newIDs {
		if err := ds.recordMatch(ctx, continent, id); err != nil {
			return err
		}
	}

	return nil
}

func (ds *Datasource) recordMatch(ctx context.Context, continent riot.Continent, id string) error {
	riotMatch, err := ds.riot.Match.GetMatch(ctx, continent, id)
	if err != nil {
		return fmt.Errorf("fetching match: %w", err)
	}

	var participants [10]Participant
	for i := range 10 {
		puuid := riotMatch.Metadata.Participants[i]
		participants[i] = NewParticipant(RiotMatchToParticipant(*riotMatch, puuid))
	}

	match := NewMatch(WithRiotMatch(riotMatch))

	err = ds.store.RecordMatch(ctx, match, participants)
	if err != nil {
		return fmt.Errorf("saving match: %w", err)
	}

	return nil
}

// ListNewMatches returns the 100 most recent matches that are not in store.
func (ds *Datasource) ListNewMatches(ctx context.Context, region riot.Region, puuid string) ([]string, error) {
	// use the continent for now
	continent := riot.RegionToContinent(region)

	ids, err := ds.riot.Match.GetMatchList(ctx, continent, puuid, "queue=420&start=0&count=100")
	if err != nil {
		return nil, fmt.Errorf("fetching ids: %w", err)
	}

	matchIDs := []string{}
	for _, id := range ids {
		matchIDs = append(matchIDs, string(id))
	}

	newIDs, err := ds.store.GetNewMatchIDs(ctx, matchIDs)
	return newIDs, err
}

// UpdateSummoner syncs a summoner-rank with riot.
func (ds *Datasource) UpdateSummoner(ctx context.Context, region riot.Region, puuid string) error {
	name, tag, err := ds.GetRiotName(ctx, puuid)
	if err != nil {
		return err
	}

	entries, err := ds.riot.League.GetLeagueEntriesByPUUID(ctx, region, puuid)
	if err != nil {
		return fmt.Errorf("fetching summoner id: %w", err)
	}

	soloq := findSoloQLeagueEntry(entries)

	var rank *RankDetail
	if soloq != nil {
		rd := NewRankDetail(WithRiotLeagueEntry(*soloq))
		rank = &rd
	}

	err = ds.store.RecordSummoner(ctx,
		Summoner{
			PUUID:      puuid,
			Name:       name,
			Tagline:    tag,
			Platform:   "",
			SummonerID: "",
		},
		RankStatus{
			PUUID:         puuid,
			EffectiveDate: time.Now(),
			Detail:        rank,
		},
	)
	if err != nil {
		return fmt.Errorf("saving summoner: %w", err)
	}

	return nil
}

// GetRank returns the up-to-date rank for a summoner.
func (ds *Datasource) GetRank(ctx context.Context, puuid string) (*RankDetail, error) {
	panic("")
}

func findSoloQLeagueEntry(entries riot.LeagueList) (soloq *riot.LeagueEntry) {
	for _, entry := range entries {
		if entry.QueueType == riot.QueueTypeRankedSolo5x5 {
			return &entry
		}
	}

	return nil
}
