package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type Datasource struct {
	riot *riot.Client

	store *Store
}

func NewDatasource(client *riot.Client, store *Store) *Datasource {
	return &Datasource{client, store}
}

func (ds *Datasource) GetLiveMatch() {

}

// GetMatchEvents returns all item events.
func (ds *Datasource) GetMatchEvents(ctx context.Context, id string) ([]ItemFrame, error) {
	events, err := ds.store.GetItemEvents(ctx, id)
	if err != nil {
		return nil, err
	}

	items := makeItemProgression(events)
	return items, nil
}

func (ds *Datasource) GetRiotName(ctx context.Context, platform, puuid string) (name, tagline string, err error) {
	riotRegion := riot.PlatformToRegion(platform)
	account, err := ds.riot.GetAccountByRiotID(ctx, riotRegion, name, tagline)
	if err != nil {
		return "", "", err
	}

	return account.GameName, account.TagLine, nil
}

func (ds *Datasource) GetPUUID(ctx context.Context, platform, name, tagline string) (puuid string, err error) {
	riotRegion := riot.PlatformToRegion(platform)
	account, err := ds.riot.GetAccountByRiotID(ctx, riotRegion, name, tagline)
	if err != nil {
		return "", err
	}

	return account.PUUID, nil
}

func (ds *Datasource) GetSummonerID(ctx context.Context, platform, puuid string) (string, error) {
	summoner, err := ds.riot.GetSummonerByPuuid(ctx, platform, puuid)
	if err != nil {
		return "", err
	}

	return summoner.Id, nil
}

func (ds *Datasource) RecordMatchTimeline(ctx context.Context, platform, id string) error {
	riotRegion := riot.PlatformToRegion(platform)
	timeline, err := ds.riot.GetMatchTimeline(ctx, riotRegion, id)
	if err != nil {
		return err
	}

	itemEvents := makeItemEvents(*timeline)
	err = ds.store.RecordItemEvents(ctx, itemEvents)
	if err != nil {
		return err
	}

	return nil
}

// UpdateMatchlist syncs the matchlist of puuid from start to start + count.
func (ds *Datasource) UpdateMatchlist(ctx context.Context, platform, puuid string, start, count int) error {
	riotRegion := riot.PlatformToRegion(platform)
	query := fmt.Sprintf("queue=420&start=%d&count=%d", start, count)
	ids, err := ds.riot.GetMatchIDsByPUUID(ctx, riotRegion, puuid, query)
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
		if err := ds.RecordMatch(ctx, platform, id); err != nil {
			return err
		}
	}

	return nil
}

// RecordMatch records a riot match.
func (ds *Datasource) RecordMatch(ctx context.Context, platform string, id string) error {
	riotRegion := riot.PlatformToRegion(platform)

	riotMatch, err := ds.riot.GetMatch(ctx, riotRegion, string(id))
	if err != nil {
		return err
	}

	var participants [10]Participant
	for i := range 10 {
		puuid := riotMatch.Metadata.Participants[i]
		participants[i] = NewParticipant(RiotMatchToParticipant(*riotMatch, puuid))
	}

	match := NewMatch(WithRiotMatch(riotMatch))

	err = ds.store.RecordMatch(ctx, match, participants)
	if err != nil {
		return err
	}

	return nil
}

func (ds *Datasource) ListNewMatches(ctx context.Context, platform string, puuid string) ([]string, error) {
	ids, err := ds.riot.GetMatchIDsByPUUID(ctx, "", puuid, "queue=420&start=0&count=100")
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

// RecordSummoner syncs the summoner's name and their rank.
func (ds *Datasource) RecordSummoner(ctx context.Context, platform, puuid string) error {
	riotRegion := riot.PlatformToRegion(platform)
	account, err := ds.riot.GetAccountByPuuid(ctx, riotRegion, puuid)
	if err != nil {
		return fmt.Errorf("fetching account: %w", err)
	}

	entries, err := ds.riot.GetLeagueEntries(ctx, platform, puuid)
	if err != nil {
		return fmt.Errorf("fetching summoner id: %w", err)
	}

	soloq := findSoloQLeagueEntry(entries)
	err = ds.store.RecordSummoner(
		ctx,
		Summoner{account.PUUID, account.GameName, account.TagLine, platform, ""},
		NewRankStatus(WithRiotLeagueEntry(puuid, time.Now(), soloq)),
	)
	if err != nil {
		return fmt.Errorf("saving summoner: %w", err)
	}

	return nil
}

func (ds *Datasource) GetChampions(ctx context.Context, puuid string) ([]SummonerChampion, error) {
	return ds.store.GetChampions(ctx, puuid)
}

func (ds *Datasource) GetMatch(ctx context.Context, id string) (Match, [10]MatchSummoner, error) {
	match, participants, err := ds.store.GetMatch(ctx, id)
	if err != nil {
		return Match{}, [10]MatchSummoner{}, err
	}

	var summoners [10]MatchSummoner
	for i := range 10 {
		summoner, err := ds.store.GetSummoner(ctx, participants[i].PUUID)
		if err != nil {
			return Match{}, [10]MatchSummoner{}, err
		}

		summoners[i] = MatchSummoner{
			Name:        summoner.Name + "#" + summoner.Tagline,
			Rank:        nil,
			Participant: participants[i],
		}
	}

	return match, summoners, nil
}

func (ds *Datasource) GetMatchlist(ctx context.Context, puuid string, page int) ([]SummonerMatch, error) {
	return ds.store.GetMatches(ctx, puuid, page)
}

func (ds *Datasource) GetRank(ctx context.Context, puuid string) (*RankDetail, error) {
	return ds.store.GetRank(ctx, puuid)
}

func (ds *Datasource) GetSummoner(ctx context.Context, puuid string) (name, platform string, err error) {
	summoner, err := ds.store.GetSummoner(ctx, puuid)
	if err != nil {
		return "", "", err
	}
	return summoner.Name + "#" + summoner.Tagline, platform, nil
}

func (ds *Datasource) SearchSummoner(ctx context.Context, q string) ([]SearchResult, error) {
	return ds.store.SearchSummoner(ctx, q)
}

func findSoloQLeagueEntry(entries []*riot.LeagueEntry) (soloq *riot.LeagueEntry) {
	for _, entry := range entries {
		if entry.QueueType == riot.QueueTypeRankedSolo5x5 {
			return entry
		}
	}
	return nil
}
