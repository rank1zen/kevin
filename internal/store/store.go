package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
)

// Should be unexported
func CreateListRankOption(history []postgres.Match) [][2]postgres.ListRankOption {
	opt := [][2]postgres.ListRankOption{}

	for i, m := range history {
		rankBefore := postgres.ListRankOption{
			Offset: 0,
			Limit:  100,
			Recent: true,
		}

		rankAfter := postgres.ListRankOption{
			Offset: 0,
			Limit:  100,
			Recent: true,
		}

		if i > 0 {
			rankBefore.Start = &history[i-1].Date
		}

		rankBefore.End = &m.Date

		if i < len(history)-1 {
			rankAfter.End = &history[i+1].Date
		}

		rankAfter.Start = &m.Date

		opt = append(opt, [2]postgres.ListRankOption{rankBefore, rankAfter})
	}

	return opt
}

// Store manages connections with a postgres database.
type Store struct {
	Pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *internal.Store {
	return &internal.Store{
		Profile:         &ProfileStore{Pool: pool},
		Match:           &MatchStore{Pool: pool},
		SummonerStats:   &SummonerStatsStore{Pool: pool},
		LiveMatch:       &LiveMatchStore{Pool: pool},
		LiveMatchStatus: &LiveMatchStatusStore{Pool: pool},
	}
}

// chooseStatusID chooses some id that is suitable.
func chooseStatusID(statusIDs []int) *int {
	if len(statusIDs) == 0 {
		return nil
	}
	return &statusIDs[0]
}

// getRank returns both status and detail. It will error if both are not found.
func (db *Store) getRank(ctx context.Context, statusID int) (postgres.RankStatus, *postgres.RankDetail, error) {
	rankStore := postgres.RankStore{Tx: db.Pool}

	status, err := rankStore.GetRankStatus(ctx, statusID)
	if err != nil {
		return postgres.RankStatus{}, nil, err
	}

	detail, err := rankStore.GetRankDetail(ctx, statusID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return status, nil, nil
		}

		return postgres.RankStatus{}, nil, err
	}

	return status, &detail, err
}

func (db *Store) getMostRecentRank(ctx context.Context, puuid riot.PUUID) (m postgres.RankStatus, n *postgres.RankDetail, err error) {
	rankStore := postgres.RankStore{Tx: db.Pool}

	ids, err := rankStore.ListRankIDs(ctx, puuid, postgres.ListRankOption{Limit: 1, Recent: true})
	if err != nil {
		return m, n, err
	}

	if len(ids) != 1 {
		return m, n, errors.New("ListRankIDS did not return exactly one id")
	}

	id := ids[0]

	return db.getRank(ctx, id)
}
