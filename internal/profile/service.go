package profile

import (
	"context"
	"errors"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type ProfileService struct {
	riot  *riot.Client
	store ProfileStore
}

// GetProfileRequest represents the request payload for retrieving a summoner profile.
type GetProfileRequest struct {
	Region string `json:"region"`
	Name   string `json:"name"`
	Tag    string `json:"tag"`
}

// GetProfile retrieves a summoner's profile. NOTE: currently always updates the
// profile.
func (s *ProfileService) GetProfile(ctx context.Context, req GetProfileRequest) (*Profile, error) {
	if req.Region == "" {
		req.Region = "NA1"
	}

	riotRegion := riot.Region(req.Region)

	account, err := s.riot.Account.GetAccountByRiotID(ctx, riotRegion, req.Name, req.Tag)
	if err != nil {
		return nil, err
	}

	entries, err := s.riot.League.GetLeagueEntriesByPUUID(ctx, riotRegion, account.PUUID.String())
	if err != nil {
		return nil, err
	}

	var soloq *riot.LeagueEntry
	if result, err := findSoloQLeagueEntry(entries); err == nil {
		soloq = result
	}

	profile := mapProfile(account, soloq, time.Now().In(time.UTC))

	if err = s.store.RecordProfile(ctx, profile); err != nil {
		return nil, err
	}

	storeProfile, err := s.store.GetProfile(ctx, string(account.PUUID))
	if err != nil {
		return nil, err
	}

	return storeProfile, nil
}

// UpdateProfileRequest represents the request payload for updating a summoner profile.
type UpdateProfileRequest struct {
	Region string `json:"region"`
	Name   string `json:"name"`
	Tag    string `json:"tag"`
}

// UpdateProfile updates a summoner's profile.
func (s *ProfileService) UpdateProfile(ctx context.Context, req UpdateProfileRequest) error {
	if req.Region == "" {
		req.Region = "NA1"
	}

	riotRegion := riot.Region(req.Region)

	account, err := s.riot.Account.GetAccountByRiotID(ctx, riotRegion, req.Name, req.Tag)
	if err != nil {
		return err
	}

	entries, err := s.riot.League.GetLeagueEntriesByPUUID(ctx, riotRegion, account.PUUID.String())
	if err != nil {
		return err
	}

	var soloq *riot.LeagueEntry
	if result, err := findSoloQLeagueEntry(entries); err == nil {
		soloq = result
	}

	profile := mapProfile(account, soloq, time.Now().In(time.UTC))

	if err = s.store.RecordProfile(ctx, profile); err != nil {
		return err
	}

	return nil
}

// GetRankHistoryRequest represents the request payload for retrieving a summoner's rank history.
type GetRankHistoryRequest struct {
	Region string `json:"region"`
	Name   string `json:"name"`
	Tag    string `json:"tag"`
}

// GetRankHistory retrieves a summoner's rank history.
func (s *ProfileService) GetRankHistory(ctx context.Context, req GetRankHistoryRequest) {
	panic("not implemented")
}

type GetSummonerChampionsRequest struct {
	Region  string     `json:"region"`
	PUUID   string     `json:"puuid"`
	StartTS *time.Time `json:"startTs"`
	EndTS   *time.Time `json:"endTs"`
}

func (s *ProfileService) GetSummonerChampions(ctx context.Context, req GetSummonerChampionsRequest) ([]ChampionAverage, error) {
	champions, err := s.store.GetChampions(ctx, req.PUUID, *req.StartTS, *req.EndTS)
	if err != nil {
		return nil, err
	}

	return champions, nil
}

func (s *ProfileService) GetMatchlist(ctx context.Context, req GetProfileRequest) (*Profile, error) {
	panic("not implemented")
}

func (s *ProfileService) GetChampionStats(ctx context.Context, req GetProfileRequest) (*Profile, error) {
	panic("not implemented")
}

// findSoloQLeagueEntry is a helper function to find the solo queue entry from a
// list of league entries. It will return an error if not found.
func findSoloQLeagueEntry(entries []riot.LeagueEntry) (*riot.LeagueEntry, error) {
	for _, entry := range entries {
		if entry.QueueType == "RANKED_SOLO_5x5" {
			return &entry, nil
		}
	}

	return nil, errors.New("solo queue entry not found")
}

func mapProfile(account *riot.Account, rank *riot.LeagueEntry, effectiveDate time.Time) *Profile {
	result := Profile{
		PUUID:   string(account.PUUID),
		Name:    account.GameName,
		Tagline: account.TagLine,
		Rank: RankStatus{
			PUUID:         string(account.PUUID),
			EffectiveDate: effectiveDate,
			Detail:        nil,
		},
	}

	if rank != nil {
		result.Rank.Detail = &RankDetail{
			Wins:     rank.Wins,
			Losses:   rank.Losses,
			Tier:     string(rank.Tier),
			Division: string(rank.Division),
			LP:       rank.LeaguePoints,
		}
	}

	return &result
}
