package internal

import (
	"errors"
)

var (
	// ErrUnknownStoreError is a generic store error.
	ErrUnknownStoreError = errors.New("unknown store error")

	// ErrSummonerNotFound is returned by Store.GetSummoner when a summoner
	// is not found in the store.
	ErrSummonerNotFound = errors.New("summoner not found")

	// ErrRankUnavailable indicates that the store does not have an
	// available rank record for a summoner. Returned by GetRank.
	ErrRankUnavailable = errors.New("rank unavailable")

	// ErrMatchNotFound is returned by Store.GetMatch when a match is not
	// found in store.
	// TODO: rename to ErrMatchNotInStore.
	ErrMatchNotFound = errors.New("match not found")

	ErrMatchMissingParticipants = errors.New("match missing participants")
)
