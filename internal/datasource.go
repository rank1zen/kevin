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
// Region parameters specify the region to search. TODO: region parameters not
// implemented.
//
// TODO: Datasource should be able to decide when to call the riot API, and
// when to use cache.
type Datasource struct {
	// probably want to cache something

	riot *riot.Client

	store Store
}

func NewDatasource(client *riot.Client, store Store) *Datasource {
	return &Datasource{client, store}
}

// GetMatchDetail will fetch riot if the match is not in store.
func (ds *Datasource) GetMatchDetail(ctx context.Context, region riot.Region, matchID string) (MatchDetail, error) {
	z := MatchDetail{}

	detail, err := ds.store.GetMatchDetail(ctx, matchID)
	if err != nil {
		if errors.Is(err, ErrMatchNotFound) {
			riotMatch, err := ds.riot.Match.GetMatch(ctx, region, matchID)
			if err != nil {
				return z, err
			}

			mapper := RiotToMatchMapper{
				Match: *riotMatch,
			}

			match := mapper.Map()
			err = ds.store.RecordMatch(ctx, match)
			if err != nil {
				return z, err
			}
		} else {
			return z, err
		}
	}

	detail, err = ds.store.GetMatchDetail(ctx, matchID)
	if err != nil {
		return z, err
	}

	return detail, nil
}

// TODO: rename to SearchProfile
func (ds *Datasource) Search(ctx context.Context, region riot.Region, q string) ([]SearchResult, error) {
	results, err := ds.store.SearchSummoner(ctx, q)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// NOTE: never fetches riot.
func (ds *Datasource) GetSummonerChampions(ctx context.Context, region riot.Region, puuid riot.PUUID, start, end time.Time) ([]SummonerChampion, error) {
	champions, err := ds.store.GetChampions(ctx, puuid, start, end)
	if err != nil {
		return nil, err
	}

	return champions, nil
}

// NOTE: it always fetches the riot api (for now), so it is always accurate.
func (ds *Datasource) GetMatchHistory(ctx context.Context, region riot.Region, puuid riot.PUUID, start, end time.Time) ([]SummonerMatch, error) {
	options := soloQMatchFilter(start, end)
	ids, err := ds.riot.Match.GetMatchList(ctx, region, string(puuid), options)
	if err != nil {
		return nil, fmt.Errorf("fetching ids: %w", err)
	}

	matchIDs := []string{}
	for _, id := range ids {
		matchIDs = append(matchIDs, id)
	}

	newIDs, err := ds.store.GetNewMatchIDs(ctx, matchIDs)
	if err != nil {
		return nil, err
	}

	// TODO: put these in batch
	for _, id := range newIDs {
		riotMatch, err := ds.riot.Match.GetMatch(ctx, region, id)
		if err != nil {
			return nil, err
		}

		match := RiotToMatchMapper{Match: *riotMatch}.Map()

		err = ds.store.RecordMatch(ctx, match)
		if err != nil {
			return nil, err
		}
	}

	final, err := ds.store.GetMatchHistory(ctx, puuid, start, end)
	if err != nil {
		return nil, err
	}

	return final, nil
}

// GetLiveMatch fetches riot for the a live match. If summoner is not in
// a game, return ErrNoLiveMatch.
func (ds *Datasource) GetLiveMatch(ctx context.Context, region riot.Region, puuid riot.PUUID) (LiveMatch, error) {
	riotGame, err := ds.riot.Spectator.GetLiveMatch(ctx, region, string(puuid))
	if err != nil {
		return LiveMatch{}, err
	}

	mapper := RiotToLiveMatchMapper{Match: *riotGame}

	return mapper.Map(), nil
}

// GetRiotName returns the Riot ID (name#tag) associated with puuid.
func (ds *Datasource) GetRiotName(ctx context.Context, puuid riot.PUUID) (name, tag string, err error) {

	// Using NA for now since puuid is globally unique ...
	account, err := ds.riot.Account.GetAccountByPUUID(ctx, riot.RegionNA1, puuid.String())
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
	// Using NA for now...
	account, err := ds.riot.Account.GetAccountByRiotID(ctx, riot.RegionNA1, name, tag)
	if err != nil {
		if errors.Is(err, riot.ErrNotFound) {
			return "", ErrSummonerDoesNotExist
		}

		return "", err
	}

	return riot.PUUID(account.PUUID), nil
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

// GetProfileDetail returns a profile.
//
// NOTE: it always fetches the riot api (for now), so it is always accurate.
func (ds *Datasource) GetProfileDetail(ctx context.Context, region riot.Region, puuid riot.PUUID) (ProfileDetail, error) {
	account, err := ds.riot.Account.GetAccountByPUUID(ctx, region, puuid.String())
	if err != nil {
		return ProfileDetail{}, err
	}

	entries, err := ds.riot.League.GetLeagueEntriesByPUUID(ctx, region, puuid.String())
	if err != nil {
		return ProfileDetail{}, err
	}

	soloq := findSoloQLeagueEntry(entries)

	s := RiotToProfileMapper{
		Account: *account,
		Rank:    soloq,
	}

	err = ds.store.RecordProfile(ctx, s.Convert())
	if err != nil {
		return ProfileDetail{}, err
	}

	profile, err := ds.store.GetProfileDetail(ctx, puuid)
	if err != nil {
		return ProfileDetail{}, err
	}

	return profile, nil
}

func (ds *Datasource) GetProfileDetailByRiotID(ctx context.Context, region riot.Region, name, tag string) (ProfileDetail, error) {
	account, err := ds.riot.Account.GetAccountByRiotID(ctx, region, name, tag)
	if err != nil {
		return ProfileDetail{}, err
	}

	entries, err := ds.riot.League.GetLeagueEntriesByPUUID(ctx, region, account.PUUID.String())
	if err != nil {
		return ProfileDetail{}, err
	}

	soloq := findSoloQLeagueEntry(entries)

	s := RiotToProfileMapper{
		Account: *account,
		Rank:    soloq,
	}

	err = ds.store.RecordProfile(ctx, s.Convert())
	if err != nil {
		return ProfileDetail{}, err
	}

	profile, err := ds.store.GetProfileDetail(ctx, account.PUUID)
	if err != nil {
		return ProfileDetail{}, err
	}

	return profile, nil
}

func (ds *Datasource) UpdateProfile(ctx context.Context, region riot.Region, puuid riot.PUUID) error {
	account, err := ds.riot.Account.GetAccountByPUUID(ctx, region, puuid.String())
	if err != nil {
		return err
	}

	entries, err := ds.riot.League.GetLeagueEntriesByPUUID(ctx, region, puuid.String())
	if err != nil {
		return err
	}

	soloq := findSoloQLeagueEntry(entries)

	s := RiotToProfileMapper{
		Account: *account,
		Rank:    soloq,
	}

	err = ds.store.RecordProfile(ctx, s.Convert())
	if err != nil {
		return err
	}

	return nil
}

func (ds *Datasource) UpdateProfileByRiotID(ctx context.Context, region riot.Region, name, tag string) error {
	account, err := ds.riot.Account.GetAccountByRiotID(ctx, region, name, tag)
	if err != nil {
		return err
	}

	entries, err := ds.riot.League.GetLeagueEntriesByPUUID(ctx, region, account.PUUID.String())
	if err != nil {
		return err
	}

	soloq := findSoloQLeagueEntry(entries)

	s := RiotToProfileMapper{
		Account: *account,
		Rank:    soloq,
	}

	err = ds.store.RecordProfile(ctx, s.Convert())
	if err != nil {
		return err
	}

	return nil
}

func findSoloQLeagueEntry(entries riot.LeagueList) (soloq *riot.LeagueEntry) {
	for _, entry := range entries {
		if entry.QueueType == riot.QueueTypeRankedSolo5x5 {
			return &entry
		}
	}

	return nil
}

func soloQMatchFilter(start, end time.Time) riot.MatchListOptions {
	options := riot.MatchListOptions{
		StartTime: new(int64),
		EndTime:   new(int64),
		Queue:     new(int),
		Type:      nil,
		Start:     0,
		Count:     100,
	}

	*options.Queue = 420
	*options.StartTime = start.Unix()
	*options.EndTime = end.Unix()

	return options
}
