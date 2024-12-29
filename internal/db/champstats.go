package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal"
)

func (db *DB) GetChampionStatList(ctx context.Context, puuid internal.PUUID, season internal.Season) ([]internal.ChampionStats, error) {
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

			avg(mp.kills)                                                       AS kda_kills,
			avg(mp.deaths)                                                      AS kda_deaths,
			avg(mp.assists)                                                     AS kda_assists,
			avg(round(mp.kills/team_stats.kills, 2))                            AS kda_participation,
			avg(mp.creep_score)                                                 AS cs_raw,
			avg(mp.creep_score)                                                 AS cs_per10,
			avg(mp.total_damage_dealt_to_champions)                             AS dmg_raw,
			avg(round(mp.total_damage_dealt_to_champions/team_total.damage, 2)) AS dmg_percentage_team,
			avg(mp.gold_earned)                                                 AS gold_raw,
			avg()                                                               AS gold_percentage_team,
			avg(mp.vision_score)                                                AS vis_raw
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
	)
	`, puuid, season)

	stats, err := pgx.CollectRows(rows, pgx.RowToStructByPos[ProfileChampionStat])
	if err != nil {
		return ProfileChampionStatList{}, err
	}

	return ProfileChampionStatList{
		Puuid:  puuid,
		Season: season,
		List:   stats,
	}, nil
}
