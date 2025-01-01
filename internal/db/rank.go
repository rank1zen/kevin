package db

import (
	"context"

	"github.com/rank1zen/yujin/internal"
)

func (db *DB) GetRankList(ctx context.Context, puuid internal.PUUID) ([]internal.RankRecord, error) {
	panic("not implemented")
	// rows, err := db.pool.Query(ctx, `
	// SELECT
	// 	wins,
	// 	losses,
	// 	division,
	// 	tier,
	// 	league_points AS lp,
	// 	entered_at    AS timestamp
	// FROM
	// 	league_records
	// WHERE
	// 	summoner_id = $1
	// ORDER BY
	// 	timestamp DESC;
	// `, puuid)
	//
	// return pgx.CollectRows(rows, pgx.RowToStructByName[RankRecord])
}
