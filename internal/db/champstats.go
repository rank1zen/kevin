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
	SELECT
		puuid,
		champion_id,
		count(*),
		0.6 -- &s.WinPercentage,
		sum(win = true),
		sum(win = false),
		10 -- &s.LpDelta,
		avg(participant.kils),
		avg(participant.deaths),
		avg(participant.assists),
		avg((participant.kills + participant.assists) / team.kills)
		avg(participant.creep_score),
		avg(participant.creep_score * 60 / extract(epoch from matches.duration)),
		avg(participant.total_damage_dealt_to_champions)
		avg(participant.total_damage_dealt_to_champions / team.total_damage_dealt_to_champions)
		0.1 -- &s.DamageDelta,
		avg(participant.gold_earned),
		avg(participant.gold_earned / team.gold_earned)
		0.1 -- &s.GoldDelta,
		avg(participant.vision_score)
	FROM
		match_participants AS participant
	JOIN
		match_team_stats AS team USING (match_id, team_id)
	JOIN
		matches USING (match_id)
	GROUP BY
		champion_id;
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
