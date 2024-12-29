package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/pgxutil"
)

func createLeagueRecord(ctx context.Context, conn pgxutil.Conn, profile internal.Profile) error {
	row := pgx.NamedArgs{}

	row["summoner_id"] = profile.SummonerID
	row["entered_at"] = profile.RecordDate

	if profile.Rank == nil {
		row["is_ranked"] = false
	} else {
		row["is_ranked"] = true
		row["league_id"] = profile.Rank.LeagueID
		row["tier"] = profile.Rank.Tier
		row["division"] = profile.Rank.Division
		row["league_points"] = profile.Rank.LP
		row["wins"] = profile.Rank.Wins
		row["losses"] = profile.Rank.Losses
	}

	_, err := conn.Exec(ctx, `
	INSERT INTO league_records (
		summoner_id,
		valid_from,
		valid_to,
		entered_at,
		is_ranked,
		league_id,
		tier,
		division,
		league_points,
		wins,
		losses
	)
	VALUES (
		@summoner_id,
		@valid_from,
		@valid_to,
		@entered_at,
		@is_ranked,
		@league_id,
		@tier,
		@division,
		@league_points,
		@wins,
		@losses
	)
	`, row)

	return err
}
