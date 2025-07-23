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

// Datasource manages interaction between the riot API and an internal store.
//
// Region parameters specify the region to search.
type Datasource struct {
	// probably want to cache something

	riot *riot.Client

	store Store
}

func NewDatasource(client *riot.Client, store Store) *Datasource {
	return &Datasource{client, store}
}

func (ds *Datasource) GetStore() Store {
	return ds.store
}

// ZUpdateMatchHistory fetches every match played from start to end.
func (ds *Datasource) ZUpdateMatchHistory(ctx context.Context, region riot.Region, puuid riot.PUUID, start, end time.Time) error {
	options := riot.MatchListOptions{
		StartTime: new(int64),
		EndTime:   new(int64),
		Queue:     new(int),
		Start:     0,
		Count:     100,
	}

	*options.Queue = 420
	*options.StartTime = start.Unix()
	*options.EndTime = end.Unix()

	ids, err := ds.riot.Match.GetMatchList(ctx, region, string(puuid), options)
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
		if err := ds.recordMatch(ctx, region, id); err != nil {
			return err
		}
	}

	return nil
}

// GetLiveMatch fetches riot for the a live match. If summoner is not in
// a game, return ErrNoLiveMatch.
func (ds *Datasource) GetLiveMatch(ctx context.Context, region riot.Region, puuid riot.PUUID) (LiveMatch, error) {
	riotGame, err := ds.riot.Spectator.GetLiveMatch(ctx, region, string(puuid))
	if err != nil {
		return LiveMatch{}, err
	}

	return NewLiveMatch(WithRiotLiveMatch(*riotGame)), nil
}

// GetRiotName returns the Riot ID (name#tag) associated with puuid.
func (ds *Datasource) GetRiotName(ctx context.Context, puuid riot.PUUID) (name, tag string, err error) {

	// Using AMER for now since puuid is globally unique ...
	account, err := ds.riot.Account.GetAccountByPUUID(ctx, riot.ContinentAmericas, string(puuid))
	if err != nil {
		return "", "", err
	}

	return account.GameName, account.TagLine, nil
}

// GetPUUID fetches riot for the puuid matching name#tag if it exists,
// otherwise, it returns [ErrSummonerDoesNotExist].
//
// Note that internal store might be stale.
func (ds *Datasource) GetPUUID(ctx context.Context, name, tag string) (riot.PUUID, error) {
	// Using AMER for now...
	account, err := ds.riot.Account.GetAccountByRiotID(ctx, riot.ContinentAmericas, name, tag)
	if err != nil {
		if errors.Is(err, riot.ErrNotFound) {
			return "", ErrSummonerDoesNotExist
		}

		return "", err
	}

	return riot.PUUID(account.PUUID), nil
}

// UpdateMatchHistory syncs the matchlist of summoner with puuid from start to
// start + count.
//
// Deprecated: We switched to using date instead of index for match list.
func (ds *Datasource) UpdateMatchHistory(ctx context.Context, region riot.Region, puuid riot.PUUID, start, count int) error {
	options := riot.MatchListOptions{
		Queue:     new(int),
		Start:     start,
		Count:     count,
	}

	*options.Queue = 420

	ids, err := ds.riot.Match.GetMatchList(ctx, region, string(puuid), options)
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
		if err := ds.recordMatch(ctx, region, id); err != nil {
			return err
		}
	}

	return nil
}

func (ds *Datasource) recordMatch(ctx context.Context, region riot.Region, id string) error {
	riotMatch, err := ds.riot.Match.GetMatch(ctx, region, id)
	if err != nil {
		return fmt.Errorf("fetching match: %w", err)
	}

	err = ds.store.RecordMatch(ctx, NewMatch(WithRiotMatch(riotMatch)))
	if err != nil {
		return fmt.Errorf("saving match: %w", err)
	}

	return nil
}

// ListNewMatches returns the 100 most recent matches that are not in store.
func (ds *Datasource) ListNewMatches(ctx context.Context, region riot.Region, puuid riot.PUUID) ([]string, error) {
	ids, err := ds.riot.Match.GetMatchList(ctx, region, string(puuid), riot.MatchListOptions{Count: 100})
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
func (ds *Datasource) UpdateSummoner(ctx context.Context, region riot.Region, puuid riot.PUUID) error {
	name, tag, err := ds.GetRiotName(ctx, puuid)
	if err != nil {
		return err
	}

	entries, err := ds.riot.League.GetLeagueEntriesByPUUID(ctx, region, string(puuid))
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
func (ds *Datasource) GetRank(ctx context.Context, puuid riot.PUUID) (*RankDetail, error) {
	return nil, nil
}

func findSoloQLeagueEntry(entries riot.LeagueList) (soloq *riot.LeagueEntry) {
	for _, entry := range entries {
		if entry.QueueType == riot.QueueTypeRankedSolo5x5 {
			return &entry
		}
	}

	return nil
}
