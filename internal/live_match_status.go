package internal

import (
	"context"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

// LiveMatchStatusStore handles live match status records. Live match status
// indicates the status of a live match.
type LiveMatchStatusStore interface {
	GetLiveMatchStatus(ctx context.Context, region riot.Region, id string) (*LiveMatchStatus, error)

	CreateLiveMatchStatus(ctx context.Context, status *LiveMatchStatus) error

	ExpireLiveMatch(ctx context.Context, region riot.Region, id string) error
}

type LiveMatchStatus struct {
	Region riot.Region

	ID string

	// Date is the datetime of the start of the match.
	Date time.Time

	// Expired is true if the match has ended.
	Expired bool
}
