package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/kevin/internal/riot"
)

type LiveMatchStatusObjects struct{ Tx Tx }

func (db *LiveMatchStatusObjects) Get(ctx context.Context, id string) (*LiveMatchStatus, error) {
	rows, err := db.Tx.Query(ctx, `
		SELECT * FROM live_match_status WHERE match_id = $1
	`, id)
	if err != nil {
		return nil, err
	}

	champions, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[LiveMatchStatus])
	if err != nil {
		return nil, err
	}

	return champions, nil
}

func (db *LiveMatchStatusObjects) Create(ctx context.Context, status *LiveMatchStatus) error {
	_, err := db.Tx.Exec(ctx, `
		INSERT INTO live_match_status (region, match_id, date, expired)
		VALUES ($1, $2, $3, $4)
	`, status.Region, status.ID, status.Date, status.Expired)
	return err
}

func (db *LiveMatchStatusObjects) Update(ctx context.Context, id string, update LiveMatchStatusUpdate) error {
	_, err := db.Tx.Exec(ctx, `
		UPDATE live_match_status SET expired = $1 WHERE match_id = $2;
	`, update.Expired, id)
	return err
}

type LiveMatchStatus struct {
	Region  riot.Region `db:"region"`
	ID      string      `db:"match_id"`
	Date    time.Time   `db:"date"`
	Expired bool        `db:"expired"`
}

type LiveMatchStatusUpdate struct {
	Expired *bool `db:"expired"`
}
