package match

import "context"

type MatchStore interface {
	CheckMatch(ctx context.Context, id string) (bool, error)

	RecordMatch(ctx context.Context, match *Match) error

	GetMatchDetail(ctx context.Context, id string) (*MatchDetail, error)
}
