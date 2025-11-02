package internal

import (
	"context"
	"errors"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type Profile struct {
	PUUID riot.PUUID

	Name, Tagline string

	Rank RankStatus
}

type SearchResult struct {
	PUUID riot.PUUID

	Name, Tagline string

	// Rank is the summoners most recent rank record in store.
	Rank *RankStatus
}

type ProfileStore interface {
	RecordProfile(ctx context.Context, profile *Profile) error

	GetProfile(ctx context.Context, puuid riot.PUUID) (*Profile, error)

	SearchSummoner(ctx context.Context, q string) ([]SearchResult, error)
}

var (
	// ErrSummonerNotFound indicates the summoner associated with some
	// puuid or name#tagline does not exist.
	ErrSummonerDoesNotExist = errors.New("summoner does not exist")

	// ErrNoLiveMatch indicates a summoner is not in a game.
	ErrNoLiveMatch = errors.New("no live match")
)

type ProfileService Datasource

type GetProfileRequest struct {
	Region *riot.Region `json:"region"`
	Name   string       `json:"name"`
	Tag    string       `json:"tag"`
}

func (s *ProfileService) GetProfile(ctx context.Context, req GetProfileRequest) (*Profile, error) {
	if req.Region == nil {
		req.Region = new(riot.Region)
		*req.Region = riot.RegionNA1
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

	profile := Profile{
		PUUID:   account.PUUID,
		Name:    account.GameName,
		Tagline: account.TagLine,
		Rank: RankStatus{
			PUUID:         account.PUUID,
			EffectiveDate: time.Time{},
			Detail: &RankDetail{
				Wins:   soloq.Wins,
				Losses: soloq.Losses,
				Rank: Rank{
					Tier:     "",
					Division: "",
					LP:       0,
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
