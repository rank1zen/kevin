package service

import (
	"context"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
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

	soloq := findSoloQLeagueEntry(entries)

	profile := internal.Profile{
		PUUID:   account.PUUID,
		Name:    account.GameName,
		Tagline: account.TagLine,
		Rank: internal.RankStatus{
			PUUID:         account.PUUID,
			EffectiveDate: time.Now().In(time.UTC),
			Detail: &internal.RankDetail{
				Wins:   soloq.Wins,
				Losses: soloq.Losses,
				Rank: internal.Rank{
					Tier:     soloq.Tier,
					Division: soloq.Division,
					LP:       soloq.LeaguePoints,
				},
			},
		},
	}

	if err = s.profile.RecordProfile(ctx, &profile); err != nil {
		return nil, err
	}

	storeProfile, err := s.profile.GetProfile(ctx, account.PUUID)
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
		*req.Region = riot.RegionNA1 // Default to NA1
	}

	account, err := s.riot.Account.GetAccountByRiotID(ctx, *req.Region, req.Name, req.Tag)
	if err != nil {
		return err
	}

	entries, err := s.riot.League.GetLeagueEntriesByPUUID(ctx, *req.Region, account.PUUID.String())
	if err != nil {
		return err
	}

	soloq := findSoloQLeagueEntry(entries)

	profile := internal.Profile{
		PUUID:   account.PUUID,
		Name:    account.GameName,
		Tagline: account.TagLine,
		Rank: internal.RankStatus{
			PUUID:         account.PUUID,
			EffectiveDate: time.Time{}, // Original had time.Time{} here
			Detail: &internal.RankDetail{
				Wins:   soloq.Wins,
				Losses: soloq.Losses,
				Rank: internal.Rank{
					Tier:     "", // Original had "" here
					Division: "", // Original had "" here
					LP:       0,  // Original had 0 here
				},
			},
		},
	}

	if err = s.profile.RecordProfile(ctx, &profile); err != nil {
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

	soloq := findSoloQLeagueEntry(entries)

	profile := internal.Profile{
		PUUID:   account.PUUID,
		Name:    account.GameName,
		Tagline: account.TagLine,
		Rank: internal.RankStatus{
			PUUID:         account.PUUID,
			EffectiveDate: time.Now().In(time.UTC),
			Detail: &internal.RankDetail{
				Wins:   soloq.Wins,
				Losses: soloq.Losses,
				Rank: internal.Rank{
					Tier:     soloq.Tier,
					Division: soloq.Division,
					LP:       soloq.LeaguePoints,
				},
			},
		},
	}

	if err = s.profile.RecordProfile(ctx, &profile); err != nil {
		return nil, err
	}

	storeProfile, err := s.profile.GetProfile(ctx, account.PUUID)
	if err != nil {
		return nil, err
	}

	return storeProfile, nil
}

type GetSummonerChampionsRequest struct {
	Region  *riot.Region `json:"region"`
	PUUID   riot.PUUID   `json:"puuid"`
	StartTS *time.Time   `json:"startTs"`
	EndTS   *time.Time   `json:"endTs"`
}

func (s *ProfileService) GetSummonerChampions(ctx context.Context, req GetSummonerChampionsRequest) ([]internal.SummonerChampion, error) {
	champions, err := s.match.GetChampions(ctx, req.PUUID, *req.StartTS, *req.EndTS)
	if err != nil {
		return nil, err
	}

	return champions, nil
}

// findSoloQLeagueEntry is a helper function to find the Solo Queue entry from a list of league entries.
// This function was implicitly used in the original profile.go.
func findSoloQLeagueEntry(entries []riot.LeagueEntry) riot.LeagueEntry {
	for _, entry := range entries {
		if entry.QueueType == "RANKED_SOLO_5x5" {
			return entry
		}
	}
	// Return a zero-value entry if not found. Error handling for this case might be needed upstream.
	return riot.LeagueEntry{}
}
