package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal"
)

func (db *DB) GetChampionStat(ctx context.Context, puuid internal.PUUID, season internal.Season) ([]internal.ChampionStats, error) {
	panic("deprecated: remove this")
}

// GetChampionList returns a something for a summoner's stats on a specific champion.
func (db *DB) GetChampionList(ctx context.Context, puuid internal.PUUID) (internal.ChampionStatsSeason, error) {
	rows, _ := db.pool.Query(ctx, `
	WITH
	team_total AS (
		SELECT
			team_id,
			match_id,
			sum(total_damage_dealt_to_champions) AS damage,
			sum(kills)                           AS kills,
			sum(gold_earned)                     AS gold
		FROM
			match_participants
		GROUP BY
			team_id, match_id
	)
	participant AS (
		SELECT
			puuid,
			champion_id,
			count(*),
			sum(lp_delta),
			avg(mp.kills),
			avg(mp.deaths),
			avg(mp.assists),
			avg(round(mp.kills/team_stats.kills, 2)),
			avg(mp.creep_score),
			avg(mp.creep_score),
			avg(mp.total_damage_dealt_to_champions),
			avg(round(mp.total_damage_dealt_to_champions/team_total.damage, 2)),
			avg(mp.gold_earned),
			avg(),
			avg(mp.vision_score)
		FROM
			match_participants mp
		JOIN
			team_stats
		WHERE
			puuid = $1
		GROUP BY
			champion_id
		ORDER BY
			champion_id
	);
	`, puuid)

	collectFn := func(row pgx.CollectableRow) (internal.ChampionStats, error) {
		var s internal.ChampionStats
		err := row.Scan(
			&s.Puuid,
			&s.Champion,
			&s.GamesPlayed,
			&s.WinPercentage,
			&s.Wins,
			&s.Losses,
			&s.LpDelta,
			&s.Kills,
			&s.Deaths,
			&s.Assists,
			&s.KillParticipation,
			&s.CreepScore,
			&s.CsPerMinute,
			&s.Damage,
			&s.DamagePercentage,
			&s.DamageDelta,
			&s.GoldEarned,
			&s.GoldPercentage,
			&s.GoldDelta,
			&s.VisionScore,
		)

		return s, err
	}

	stats, err := pgx.CollectRows(rows, collectFn)
	if err != nil {
		return internal.ChampionStatsSeason{}, err
	}

	return internal.ChampionStatsSeason{
		List: stats,
	}, nil
}
