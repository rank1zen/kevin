package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rank1zen/kevin/internal/riot"
)

var (
	ErrInvalidRankStatuID = errors.New("invalid rank status id")
)

type RankStore struct{ Tx Tx }

// CreateRankStatus creates a rank status and returns created id.
func (db *RankStore) CreateRankStatus(ctx context.Context, status RankStatus) (id int, err error) {
	err = db.Tx.QueryRow(ctx, `
		INSERT INTO RankStatus (
			puuid,
			effective_date,
			is_ranked
		)
		VALUES (
			@puuid,
			@effective_date,
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
	).Scan(&id)

	return id, err
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

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return ErrInvalidRankStatuID
			}
		}

		return err

	}

	return nil
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
				status.effective_date >= @start
		`
		args["start"] = option.Start
	}

	if option.End != nil {
		sql = sql + `
			AND
				status.effective_date < @end
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

	sql = sql + `
		OFFSET
			@offset
		LIMIT
			@limit
	`
	args["offset"] = option.Offset
	args["limit"] = option.Limit

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

func (db *RankStore) GetRankStatus(ctx context.Context, id int) (status RankStatus, err error) {
	defer errWrap(&err, "GetRankStatus")

	rows, err := db.Tx.Query(ctx, `
		SELECT
			puuid,
			effective_date,
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

func (db *RankStore) ListRanks(ctx context.Context, region string, opt LeaderBoardOption) ([]RankFull, error) {
	rows, err := db.Tx.Query(ctx, `
		WITH recent_rank_status AS (
			SELECT
				*
			FROM
				RankStatus AS out
			WHERE
				effective_date = (
					SELECT
						MAX(effective_date)
					FROM
						RankStatus AS inn
					WHERE
						out.puuid = inn.puuid
				)
				AND
				out.is_ranked = true
		)
		SELECT
			recent_rank_status.puuid,
			recent_rank_status.effective_date,
			RankDetail.rank_status_id,
			RankDetail.wins,
			RankDetail.losses,
			RankDetail.tier,
			RankDetail.division,
			RankDetail.lp
		FROM
			RankDetail
		JOIN
			recent_rank_status ON RankDetail.rank_status_id = recent_rank_status.rank_status_id
		ORDER BY
			tier DESC,
			division DESC,
			lp DESC
		OFFSET
			@start
		LIMIT
			@count
	`,
		pgx.NamedArgs{
			"count": opt.Count,
			"start": opt.Start,
		},
	)

	if err != nil {
		return nil, err
	}

	m, err := pgx.CollectRows(
		rows,
		func(row pgx.CollectableRow) (result RankFull, _ error) {
			result.Detail = &RankDetail{}
			row.Scan(
				&result.Status.PUUID,
				&result.Status.EffectiveDate,
				&result.Detail.RankStatusID,
				&result.Detail.Wins,
				&result.Detail.Losses,
				&result.Detail.Tier,
				&result.Detail.Division,
				&result.Detail.LP,
			)
			return result, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return m, nil
}

// RankStatus is always created when a rank request is made.
type RankStatus struct {
	PUUID string `db:"puuid"`

	EffectiveDate time.Time `db:"effective_date"`

	// IsRanked indicates that there exists a [RankDetail] for this status.
	IsRanked bool `db:"is_ranked"`
}

type RankDetail struct {
	RankStatusID int `db:"rank_status_id"`

	Wins int `db:"wins"`

	Losses int `db:"losses"`

	Tier string `db:"tier"`

	Division string `db:"division"`

	LP int `db:"lp"`
}

type RankFull struct {
	Status RankStatus
	Detail *RankDetail
}

type ListRankOption struct {
	// Start indicates an inclusive lower bound on the date.
	Start *time.Time

	// End indicates an exclusive upper bound on the date.
	End *time.Time

	Offset, Limit uint

	Recent bool
}

type LeaderBoardOption struct {
	Start int
	Count int
}
