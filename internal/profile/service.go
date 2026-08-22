package profile

import (
	"github.com/rank1zen/kevin/internal/riot"
)

type Service struct {
	riot  *riot.Client
	store Store
}

func NewProfileService(riot *riot.Client, store Store) *Service {
	return &Service{
		riot:  riot,
		store: store,
	}
}
