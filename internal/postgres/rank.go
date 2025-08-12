package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

type Tx interface {
	Exec(ctx context.Context, sql string, args ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type BatchTx interface {
	Queue(query string, arguments ...any) *pgx.QueuedQuery
}

type RankDetail struct {
	RankStatusID int    `db:"rank_status_id"`
	Wins         int    `db:"wins"`
	Losses       int    `db:"losses"`
	Tier         string `db:"tier"`
	Division     string `db:"division"`
	LP           int    `db:"lp"`
}

type RankStatus struct {
	RankStatusID  int       `db:"rank_status_id"`
	PUUID         string    `db:"puuid"`
	EffectiveDate time.Time `db:"effective_date"`
	EndDate       time.Time `db:"end_date"`
	IsCurrent     bool      `db:"is_current"`
	IsRanked      bool      `db:"is_ranked"`
}

func RankRecordFromPG(status RankStatus, detail *RankDetail) internal.RankRecordFrom {
	return func(r *internal.ZRankRecord) {
		r.PUUID = riot.PUUID(status.PUUID)

		r.EffectiveDate = status.EffectiveDate

		r.EndDate = &status.EndDate

		r.IsCurrent = status.IsCurrent

		if detail != nil {
			r.Detail = &internal.ZRankDetail{
				Wins:   detail.Wins,
				Losses: detail.Losses,
				Rank: internal.Rank{
					Tier:     convertStringToRiotTier(detail.Tier),
					Division: convertStringToRiotRank(detail.Division),
					LP:       detail.LP,
				},
			}
		} else {
			r.Detail = nil
		}
	}
}

type RankStore struct{ Tx Tx }

type ListRankOption struct {
	Start, End *time.Time

	Offset, Limit uint

	Recent bool
}

func (db *RankStore) CreateRankStatus(ctx context.Context, status RankStatus) error {
	_, err := db.Tx.Exec(ctx, `
		INSERT INTO RankStatus (
			puuid,
			effective_date,
			end_date,
			is_current,
			is_ranked
		)
		VALUES (
			@puuid,
			@effective_date,
			'infinity',
			true,
			@is_ranked
		)
		RETURNING
			rank_status_id;
	`,
		pgx.NamedArgs{
			"puuid":          status.PUUID,
			"effective_date": status.EffectiveDate,
			"is_ranked":      status.IsRanked,
		},
	)

	return err
}

func (db *RankStore) CreateRankDetail(ctx context.Context, detail RankDetail) error {
	_, err := db.Tx.Exec(ctx, `
		INSERT INTO RankDetail (
			rank_status_id,
			wins,
			losses,
			tier,
			division,
			lp
		)
		VALUES (
			@rank_status_id,
			@wins,
			@losses,
			@tier,
			@division,
			@lp
		);
	`,
		pgx.NamedArgs{
			"rank_status_id": detail.RankStatusID,
			"wins":           detail.Wins,
			"losses":         detail.Losses,
			"tier":           detail.Tier,
			"division":       detail.Division,
			"lp":             detail.LP,
		},
	)

	return err
}

func (db *RankStore) ListRankIDs(ctx context.Context, puuid riot.PUUID, option ListRankOption) ([]int, error) {
	var sql string
	args := pgx.NamedArgs{}

	sql = `
		SELECT
			rank_status_id
		FROM
			RankStatus status
		WHERE
			puuid = @puuid
	`

	args["puuid"] = puuid

	if option.Start != nil {
		sql = sql + `
			AND
				status.effective_date <= @start
		`
		args["start"] = option.Start
	}

	if option.End != nil {
		sql = sql + `
			AND
				status.effective_date >= @end
		`

		args["end"] = option.End
	}

	if option.Recent {
		sql = sql + `
			ORDER BY
				status.effective_date DESC
		`
	} else {
		sql = sql + `
			ORDER BY
				status.effective_date ASC
		`
	}

	args["end"] = option.End

	rows, err := db.Tx.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}

	collect := func(row pgx.CollectableRow) (id int, err error) {
		err = row.Scan(&id)
		return id, err
	}

	ids, err := pgx.CollectRows(rows, collect)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (db *RankStore) GetRankStatus(ctx context.Context, id int) (RankStatus, error) {
	rows, err := db.Tx.Query(ctx, `
		SELECT
			rank_status_id,
			puuid,
			effective_date,
			end_date,
			is_current,
			is_ranked
		FROM
			RankStatus
		WHERE
			rank_status_id = @rank_status_id;
	`,
		pgx.NamedArgs{
			"rank_status_id": id,
		},
	)
	if err != nil {
		return RankStatus{}, err
	}

	m, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[RankStatus])
	if err != nil {
		return RankStatus{}, err
	}

	return m, nil
}

func (db *RankStore) GetRankDetail(ctx context.Context, id int) (RankDetail, error) {
	rows, err := db.Tx.Query(ctx, `
		SELECT
			rank_status_id,
			wins,
			losses,
			tier,
			division,
			lp
		FROM
			RankDetail
		WHERE
			rank_status_id = @rank_status_id;
	`,
		pgx.NamedArgs{
			"rank_status_id": id,
		},
	)
	if err != nil {
		return RankDetail{}, err
	}

	m, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[RankDetail])
	if err != nil {
		return RankDetail{}, err
	}

	return m, nil
}
