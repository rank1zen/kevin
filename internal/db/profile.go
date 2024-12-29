package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/pgxutil"
)

func (db *DB) ProfileExists(ctx context.Context, puuid internal.PUUID) (bool, error) {
	var exists bool
	err := db.pool.QueryRow(ctx, `
	SELECT EXISTS (
		SELECT 1
		FROM
			profiles
		WHERE
			puuid = $1
	);`, puuid).Scan(&exists)
	return exists, err
}

func upsertProfile(ctx context.Context, conn pgxutil.Exec, profile internal.Profile) error {
	_, err := conn.Exec(ctx, `
	INSERT INTO profiles
		(puuid, name, tagline, last_updated)
	VALUES
		($1, $2, $3, $4)
	ON CONFLICT
		(puuid)
	DO UPDATE SET
		name = $2,
		tagline = $3,
		last_updated = $4;
	`, profile.Puuid, profile.Name, profile.Tagline, profile.RecordDate)

	return err
}

func (db *DB) ProfileUpdate(ctx context.Context, profile internal.Profile) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = createSummonerRecord(ctx, tx, profile)
	if err != nil {
		return fmt.Errorf("inserting summoner: %w", err)
	}

	err = createLeagueRecord(ctx, tx, profile)
	if err != nil {
		return fmt.Errorf("inserting league: %w", err)
	}

	err = upsertProfile(ctx, tx, profile)
	if err != nil {
		return fmt.Errorf("upserting profile: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

type ProfileHeader struct {
	Puuid       internal.PUUID `db:"puuid"`
	LastUpdated time.Time      `db:"last_updated"`
	RiotID      string         `db:"riot_id"`
	RiotTagLine string         `db:"riot_tagline"`
	Rank        *RankRecord    `db:"-"`
}

func (db *DB) ProfileGetHeader(ctx context.Context, puuid internal.PUUID) (ProfileHeader, error) {
	var m ProfileHeader
	err := db.pool.QueryRow(ctx, `
	SELECT
		puuid,
		last_updated,
		name,
		tagline
	FROM
		profiles
	WHERE
		puuid = $1;
	`, puuid).Scan(
		m.Puuid,
		m.LastUpdated,
		m.RiotID,
		m.RiotTagLine,
	)

	return m, err
}

type ProfileRankHistoryList []RankRecord

func (db *DB) GetRankHistory(ctx context.Context, puuid internal.PUUID) (ProfileRankHistoryList, error) {
	rows, _ := db.pool.Query(ctx, `
	SELECT
		wins,
		losses,
		division,
		tier,
		league_points AS lp,
		entered_at    AS timestamp
	FROM
		league_records
	WHERE
		summoner_id = $1
	ORDER BY
		timestamp DESC
	`, puuid)

	return pgx.CollectRows(rows, pgx.RowToStructByName[RankRecord])
}

type ProfileMatch struct {
	Puuid         internal.PUUID      `db:"puuid"`
	MatchID       internal.MatchID    `db:"match_id"`
	Patch         string              `db:"patch"`
	Date          time.Time           `db:"date"`
	Duration      time.Duration       `db:"duration"`
	Champion      internal.ChampionID `db:"champion"`
	ChampionLevel int                 `db:"champion_level"`
	Summoners     internal.SummsIDs   `db:"summoners"`
	Runes         internal.Runes      `db:"runes"`
	Items         internal.ItemIDs    `db:"items"`
	Win           bool                `db:"win"`

	Kills             int     `db:"kills"`
	Deaths            int     `db:"deaths"`
	Assists           int     `db:"assists"`
	KillParticipation float32 `db:"kill_participation"`
	CreepScore        int     `db:"creep_score"`
	CsPerMinute       float32 `db:"cs_per_minute"`
	DamageDone        int     `db:"damage_done"`
	DamagePercentage  float32 `db:"damage_percentage"`
	DamageDelta       int     `db:"damage_delta"`
	GoldEarned        int     `db:"gold_earned"`
	GoldPercentage    float32 `db:"gold_percentage"`
	GoldDelta         int     `db:"gold_delta"`
	VisionScore       int     `db:"vision_score"`
}

type ProfileMatchList struct {
	Puuid   internal.PUUID
	Page    int
	Count   int
	HasMore bool
	List    []ProfileMatch
}

func (db *DB) ProfileGetMatchList(ctx context.Context, puuid internal.PUUID, page int, ensure bool) (ProfileMatchList, error) {
	start, count := 10*page, 10

	rows, _ := db.pool.Query(ctx, `
	SELECT
		player.puuid,
		meta.match_id,
		meta.game_patch AS patch,
		meta.game_date AS date,
		meta.game_duration AS duration,
		player.champion_id AS champion,
		player.champion_level,
		stats.kills,
		stats.deaths,
		stats.assists,
		stats.kill_participation
	FROM
		match_participants AS player
	JOIN
		match_participant_stats AS stats USING (match_id, puuid)
	JOIN
		matches AS meta USING (match_id)
	WHERE
		player.puuid = $1
	ORDER BY
		game_date DESC;
	OFFSET
		$2
	LIMIT
		$3;
	`, puuid, start, count)

	matches, err := pgx.CollectRows(rows, pgx.RowToStructByName[ProfileMatch])
	if err != nil {
		return ProfileMatchList{}, err
	}

	var hasMore bool
	db.pool.QueryRow(ctx, `
	SELECT EXISTS (SELECT 1 FROM match_participants WHERE puuid = $1 OFFSET $2);
	`, puuid, start).Scan(&hasMore)

	return ProfileMatchList{
		Puuid:   puuid,
		Page:    page,
		Count:   len(matches),
		HasMore: hasMore,
		List:    matches,
	}, nil
}

type ProfileLiveGameParticipant struct {
	Puuid       internal.PUUID  `db:"puuid"`
	TeamID      internal.TeamID `db:"-"`
	RiotID      string          `db:"name"`
	RiotTagLine string          `db:"tagline"`

	Champion  internal.ChampionID `db:"-"`
	Summoners internal.SummsIDs   `db:"-"`
	Runes     internal.Runes      `db:"-"`

	Rank *RankRecord

	BannedChampion *internal.ChampionID `db:"-"`
}
