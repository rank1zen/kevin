package db

import (
	"context"
	"time"

	"github.com/rank1zen/yujin/internal/pgxutil"
	"github.com/rank1zen/yujin/internal/riotclient"
)

type RankRecord struct {
	Puuid     riotclient.PUUID `db:"-"`
	Timestamp time.Time        `db:"timestamp"`
	Wins      int              `db:"wins"`
	Losses    int              `db:"losses"`
	Tier      int              `db:"tier"`
	Division  string           `db:"division"`
	Lp        int              `db:"lp"`
}

type RankSnapshot struct {
	RankRecord

	LpDelta *int `db:"lp_delta"` // LP change after match
}

func (db *DB) RankGetRecord() (RankRecord, error) {

}

func (db *DB) RankGetSnapshot() (RankSnapshot, error) {

}

// getRankAt finds the most recent rank available for puuid at timestamp.
func getRankAt(ctx context.Context, conn pgxutil.Conn, puuid riotclient.PUUID, ts time.Time) (*RankRecord, error) {
	var record RankRecord
	record.Puuid = puuid
	err := conn.QueryRow(ctx, `
	SELECT
		entered_at AS timestamp,
		wins
		losses
		tier
		division
		league_points AS lp
	FROM
		league_records
	WHERE
		puuid = $1 AND entered_at < $2
	`, puuid, ts).Scan(
		&record.Timestamp,
		&record.Wins,
		&record.Losses,
		&record.Tier,
		&record.Division,
		&record.Lp,
	)
	if err != nil {
		return RankRecord{}, err
	}

	return record, nil
}

// findRankSnapshotAt finds the most recent rank available as well as a delta
func findRankSnapshotAt(ctx context.Context, conn pgxutil.Conn, puuid riotclient.PUUID) (RankSnapshot, error) {
	return RankSnapshot{}, nil
}

func getLatestRank() (*RankRecord, error) {

}
