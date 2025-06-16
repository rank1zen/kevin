package internal

import "errors"

var (
	ErrSummonerNotFound = errors.New("summoner not found")
	ErrNoLiveMatch = errors.New("no live match")
)
