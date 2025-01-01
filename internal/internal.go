package internal

import (
	"context"
)

type Repository interface {
	CheckProfileExists(context.Context, PUUID) (bool, error)

	UpdateProfile(context.Context, Profile) error

	GetProfile(context.Context, PUUID) (Profile, error)

	CreateMatch(context.Context, RiotMatch) error

	// TODO: remove
	GetMatchList(context.Context, PUUID, int, bool) ([]RiotMatchParticipant, error)

	// GetMatchHistory returns a paging object for a summoner's match history.
	GetMatchHistory(context.Context, PUUID) (MatchHistory, error)

	// GetMatchHistory returns the 10 participants in a match.
	GetMatch(context.Context, MatchID) ([10]Participant, error)

	// GetChampionList returns a something for a summoner's stats on a specific champion.
	GetChampionList(context.Context, PUUID) (ChampionStatsSeason, error)

	// CheckMatchIDs returns the matches not in local.
	CheckMatchIDs(context.Context, []MatchID) ([]MatchID, error)

	GetRankList(context.Context, PUUID) ([]RankRecord, error)
}

type RiotClient interface {
	GetProfile(context.Context, PUUID) (Profile, error)

	// GetLiveMatch returns the 10 current participants in the current match for PUUID
	GetLiveMatch(context.Context, PUUID) (RiotLiveMatch, error)

	// GetLiveMatch returns the 10 current participants in the current match for PUUID
	// GetLiveMatch(context.Context, PUUID) ([10]LiveParticipant, error)

	GetMatchList(context.Context, PUUID) ([]MatchID, error)

	GetMatch(context.Context, MatchID) (RiotMatch, error)
}
