package service

import "errors"

// ErrNoLiveMatch indicates a summoner is not in a game.
var ErrNoLiveMatch = errors.New("no live match")
