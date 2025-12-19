package service

import (
	"context"
	"errors"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/riotmapper"
)

// ProfileService manages profile-related operations.
type ProfileService Service

// GetProfileRequest represents the request payload for retrieving a summoner profile.
type GetProfileRequest struct {
	Region *riot.Region `json:"region"`
	Name   string       `json:"name"`
	Tag    string       `json:"tag"`
}

// GetProfile retrieves a summoner's profile.
func (s *ProfileService) GetProfile(ctx context.Context, req GetProfileRequest) (*internal.Profile, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1 // Default to NA1
	}

	account, err := s.riot.Account.GetAccountByRiotID(ctx, *req.Region, req.Name, req.Tag)
	if err != nil {
		return nil, err
	}

	entries, err := s.riot.League.GetLeagueEntriesByPUUID(ctx, *req.Region, account.PUUID.String())
	if err != nil {
		return nil, err
	}

	var soloq *riot.LeagueEntry
	if result, err := findSoloQLeagueEntry(entries); err == nil {
		soloq = result
	}

	profile := riotmapper.MapProfile(&riotmapper.Profile{
		Account:       *account,
		Rank:          soloq,
		EffectiveDate: time.Now().In(time.UTC),
	})

	if err = s.store.Profile.RecordProfile(ctx, profile); err != nil {
		return nil, err
	}

	storeProfile, err := s.store.Profile.GetProfile(ctx, account.PUUID)
	if err != nil {
		return nil, err
	}

	return storeProfile, nil
}

// UpdateProfileRequest represents the request payload for updating a summoner profile.
type UpdateProfileRequest struct {
	Region *riot.Region `json:"region"`
	Name   string       `json:"name"`
	Tag    string       `json:"tag"`
}

// UpdateProfile updates a summoner's profile.
func (s *ProfileService) UpdateProfile(ctx context.Context, req UpdateProfileRequest) error {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
	}

	account, err := s.riot.Account.GetAccountByRiotID(ctx, *req.Region, req.Name, req.Tag)
	if err != nil {
		return err
	}

	entries, err := s.riot.League.GetLeagueEntriesByPUUID(ctx, *req.Region, account.PUUID.String())
	if err != nil {
		return err
	}

	profile := internal.Profile{
		PUUID:   account.PUUID,
		Name:    account.GameName,
		Tagline: account.TagLine,
		Rank: internal.RankStatus{
			PUUID:         account.PUUID,
			EffectiveDate: time.Now().In(time.UTC),
			Detail:        nil,
		},
	}

	soloq, err := findSoloQLeagueEntry(entries)
	if err == nil {
		profile.Rank.Detail = &internal.RankDetail{
			Wins:   soloq.Wins,
			Losses: soloq.Losses,
			Rank: internal.Rank{
				Tier:     soloq.Tier,
				Division: soloq.Division,
				LP:       soloq.LeaguePoints,
			},
		}
	}

	if err = s.store.Profile.RecordProfile(ctx, &profile); err != nil {
		return err
	}

	return nil
}

// GetRankHistoryRequest represents the request payload for retrieving a summoner's rank history.
type GetRankHistoryRequest struct {
	Region *riot.Region `json:"region"`
	Name   string       `json:"name"`
	Tag    string       `json:"tag"`
}

// GetRankHistory retrieves a summoner's rank history.
func (s *ProfileService) GetRankHistory(ctx context.Context, req GetRankHistoryRequest) (*internal.Profile, error) {
	panic("not implemented")
}

type GetSummonerChampionsRequest struct {
	Region  *riot.Region `json:"region"`
	PUUID   riot.PUUID   `json:"puuid"`
	StartTS *time.Time   `json:"startTs"`
	EndTS   *time.Time   `json:"endTs"`
}

func (s *ProfileService) GetSummonerChampions(ctx context.Context, req GetSummonerChampionsRequest) ([]internal.SummonerChampion, error) {
	champions, err := s.store.SummonerStats.GetChampions(ctx, req.PUUID, *req.StartTS, *req.EndTS)
	if err != nil {
		return nil, err
	}

	return champions, nil
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
