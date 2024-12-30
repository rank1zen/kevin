package internal

import (
	"context"
)

type Repository interface {
	CheckProfileExists(context.Context, PUUID) (bool, error)

	UpdateProfile(context.Context, Profile) error

	GetProfile(context.Context, PUUID) (Profile, error)

	GetRankList(context.Context, PUUID) ([]RankRecord, error)

	GetMatchList(context.Context, PUUID, int, bool) ([]MatchParticipant, error)

	GetChampionList(context.Context, PUUID) ([]ChampionStats, error)

	CreateMatch(context.Context, Match) error
}

type RiotClient interface {
	GetProfile(context.Context, PUUID) (Profile, error)

	GetLiveMatch(context.Context, PUUID) (LiveMatch, error)

	GetMatchList(context.Context, PUUID) ([]MatchID, error)

	GetMatch(context.Context, MatchID) (Match, error)
}
