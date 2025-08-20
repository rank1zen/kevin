package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

type Match struct {
	ID string `db:"match_id"`

	Date time.Time `db:"date"`

	Duration time.Duration `db:"duration"`

	Version string `db:"version"`

	WinnerID int `db:"winner"`
}

type Participant struct {
	PUUID string `db:"puuid"`

	MatchID string `db:"match_id"`

	TeamID int `db:"team"`

	ChampionID int `db:"champion"`

	ChampionLevel int `db:"champion_level"`

	TeamPosition string `db:"position"`

	SummonerIDs [2]int `db:"summoners"`

	Runes [11]int `db:"runes"`

	Items [7]int `db:"items"`

	Kills int `db:"kills"`

	Deaths int `db:"deaths"`

	Assists int `db:"assists"`

	KillParticipation float32 `db:"kill_participation"`

	CreepScore int `db:"creep_score"`

	CreepScorePerMinute float32 `db:"creep_score_per_minute"`

	DamageDealt int `db:"damage_dealt"`

	DamageTaken int `db:"damage_taken"`

	DamageDeltaEnemy int `db:"damage_delta_enemy"`

	DamagePercentageTeam float32 `db:"damage_percentage_team"`

	GoldEarned int `db:"gold_earned"`

	GoldDeltaEnemy int `db:"gold_delta_enemy"`

	GoldPercentageTeam float32 `db:"gold_percentage_team"`

	VisionScore int `db:"vision_score"`

	PinkWardsBought int `db:"pink_wards_bought"`
}

// MatchStore manages match objects.
type MatchStore struct{ Tx Tx }

func (db *MatchStore) GetSummonerChampions(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]internal.SummonerChampion, error) {
	rows, err := db.Tx.Query(ctx, `
		SELECT
			puuid,
			count(*) AS games_played,
			sum(CASE WHEN p.team = m.winner THEN 1 ELSE 0 END) AS num_matches,
			champion,
			avg(kills),
			avg(deaths),
			avg(assists),
			avg(kill_participation),
			round(avg(creep_score)),
			avg(creep_score_per_minute),
			round(avg(gold_earned)),
			round(avg(gold_delta_enemy)),
			avg(gold_percentage_team),
			round(avg(damage_dealt)),
			round(avg(damage_delta_enemy)),
			avg(damage_percentage_team),
			round(avg(damage_taken)),
			round(avg(vision_score)),
			round(avg(pink_wards_bought))
		FROM
			Participant p
		JOIN
			Match m USING (match_id)
		WHERE
			puuid = @puuid
		AND
			m.date >= @start
		AND
			m.date <= @end
		GROUP BY
			puuid, champion
		ORDER BY
			count(*) DESC;
	`,
		pgx.NamedArgs{
			"puuid": puuid,
			"start": start,
			"end":   end,
		},
	)
	if err != nil {
		return nil, err
	}

	collect := func(row pgx.CollectableRow) (m internal.SummonerChampion, err error) {
		err = row.Scan(
			&m.PUUID,
			&m.GamesPlayed,
			&m.Wins,
			&m.Champion,
			&m.AverageKillsPerGame,
			&m.AverageDeathsPerGame,
			&m.AverageAssistsPerGame,
			&m.AverageKillParticipationPerGame,
			&m.AverageCreepScorePerGame,
			&m.AverageCreepScorePerMinutePerGame,
			&m.AverageGoldEarnedPerGame,
			&m.AverageGoldDeltaEnemyPerGame,
			&m.AverageGoldPercentagePerGame,
			&m.AverageDamageDealtPerGame,
			&m.AverageDamageDeltaEnemyPerGame,
			&m.AverageDamagePercentagePerGame,
			&m.AverageDamageTakenPerGame,
			&m.AverageVisionScorePerGame,
			&m.AveragePinkWardsBoughtPerGame,
		)

		return m, err
	}

	champions, err := pgx.CollectRows(rows, collect)
	if err != nil {
		return nil, err
	}

	return champions, nil
}

func (db *MatchStore) CreateMatchInBatch(tx BatchTx, match Match) {
	tx.Queue(`
		INSERT INTO Match (
			match_id,
			date,
			duration,
			version,
			winner
		)
		VALUES (
			@match_id,
			@date,
			@duration,
			@version,
			@winner
		);
	`,
		pgx.NamedArgs{
			"match_id": match.ID,
			"date":     match.Date,
			"duration": match.Duration,
			"version":  match.Version,
			"winner":   match.WinnerID,
		},
	)
}

func (db *MatchStore) CreateMatch(ctx context.Context, match Match) error {
	_, err := db.Tx.Exec(ctx, `
		INSERT INTO Match (
			match_id,
			date,
			duration,
			version,
			winner
		)
		VALUES (
			@match_id,
			@date,
			@duration,
			@version,
			@winner
		);
	`,
		pgx.NamedArgs{
			"match_id": match.ID,
			"date":     match.Date,
			"duration": match.Duration,
			"version":  match.Version,
			"winner":   match.WinnerID,
		},
	)

	return err
}

func (db *MatchStore) CreateParticipantInBatch(tx BatchTx, participant Participant) {
	tx.Queue(`
		INSERT INTO Participant (
			match_id,
			puuid,
			team,
			position,
			champion,
			champion_level,
			summoners,
			runes,
			items,
			kills,
			deaths,
			assists,
			kill_participation,
			creep_score,
			creep_score_per_minute,
			damage_dealt,
			damage_taken,
			damage_delta_enemy,
			damage_percentage_team,
			gold_earned,
			gold_delta_enemy,
			gold_percentage_team,
			vision_score,
			pink_wards_bought
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14,
			$15,
			$16,
			$17,
			$18,
			$19,
			$20,
			$21,
			$22,
			$23,
			$24
		);
	`,
		participant.MatchID,
		participant.PUUID,
		participant.TeamID,
		participant.TeamPosition,
		participant.ChampionID,
		participant.ChampionLevel,
		participant.SummonerIDs,
		participant.Runes,
		participant.Items,
		participant.Kills,
		participant.Deaths,
		participant.Assists,
		participant.KillParticipation,
		participant.CreepScore,
		participant.CreepScorePerMinute,
		participant.DamageDealt,
		participant.DamageTaken,
		participant.DamageDeltaEnemy,
		participant.DamagePercentageTeam,
		participant.GoldEarned,
		participant.GoldDeltaEnemy,
		participant.GoldPercentageTeam,
		participant.VisionScore,
		participant.PinkWardsBought,
	)
}

func (db *MatchStore) CreateParticipant(ctx context.Context, participant Participant) error {
	_, err := db.Tx.Exec(ctx, `
		INSERT INTO Participant (
			match_id,
			puuid,
			team,
			position,
			champion,
			champion_level,
			summoners,
			runes,
			items,
			kills,
			deaths,
			assists,
			kill_participation,
			creep_score,
			creep_score_per_minute,
			damage_dealt,
			damage_taken,
			damage_delta_enemy,
			damage_percentage_team,
			gold_earned,
			gold_delta_enemy,
			gold_percentage_team,
			vision_score,
			pink_wards_bought
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14,
			$15,
			$16,
			$17,
			$18,
			$19,
			$20,
			$21,
			$22,
			$23,
			$24
		);
	`,
		participant.MatchID,
		participant.PUUID,
		participant.TeamID,
		participant.TeamPosition,
		participant.ChampionID,
		participant.ChampionLevel,
		participant.SummonerIDs,
		participant.Runes,
		participant.Items,
		participant.Kills,
		participant.Deaths,
		participant.Assists,
		participant.KillParticipation,
		participant.CreepScore,
		participant.CreepScorePerMinute,
		participant.DamageDealt,
		participant.DamageTaken,
		participant.DamageDeltaEnemy,
		participant.DamagePercentageTeam,
		participant.GoldEarned,
		participant.GoldDeltaEnemy,
		participant.GoldPercentageTeam,
		participant.VisionScore,
		participant.PinkWardsBought,
	)

	return err
}

func (db *MatchStore) GetMatch(ctx context.Context, id string) (Match, error) {
	rows, err := db.Tx.Query(ctx, `
		SELECT
			match_id,
			date,
			duration,
			version,
			winner
		FROM
			Match
		WHERE
			match_id = @match_id;
	`,
		pgx.NamedArgs{
			"match_id": id,
		},
	)
	if err != nil {
		return Match{}, err
	}

	match, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Match])
	if err != nil {
		return Match{}, err
	}

	return match, nil
}

func (db *MatchStore) GetParticipants(ctx context.Context, id string) ([]Participant, error) {
	rows, err := db.Tx.Query(ctx, `
		SELECT
			puuid,
			match_id,
			team,
			position,
			champion,
			champion_level,
			summoners,
			runes,
			items,
			kills,
			deaths,
			assists,
			kill_participation,
			creep_score,
			creep_score_per_minute,
			gold_earned,
			gold_delta_enemy,
			gold_percentage_team,
			damage_dealt,
			damage_taken,
			damage_delta_enemy,
			damage_percentage_team,
			vision_score,
			pink_wards_bought
		FROM
			Participant
		WHERE
			match_id = @match_id;
	`,
		pgx.NamedArgs{
			"match_id": id,
		},
	)

	if err != nil {
		return nil, err
	}

	participants, err := pgx.CollectRows(rows, pgx.RowToStructByName[Participant])
	if err != nil {
		return nil, err
	}

	return participants, nil
}

func (db *MatchStore) GetParticipant(ctx context.Context, puuid riot.PUUID, matchID string) (Participant, error) {
	rows, err := db.Tx.Query(ctx, `
		SELECT
			puuid,
			match_id,
			team,
			position,
			champion,
			champion_level,
			summoners,
			runes,
			items,
			kills,
			deaths,
			assists,
			kill_participation,
			creep_score,
			creep_score_per_minute,
			gold_earned,
			gold_delta_enemy,
			gold_percentage_team,
			damage_dealt,
			damage_taken,
			damage_delta_enemy,
			damage_percentage_team,
			vision_score,
			pink_wards_bought
		FROM
			Participant
		WHERE
			match_id = @match_id
			AND puuid = @puuid;
	`,
		pgx.NamedArgs{
			"puuid":    puuid.String(),
			"match_id": matchID,
		},
	)

	if err != nil {
		return Participant{}, err
	}

	participants, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Participant])
	if err != nil {
		return Participant{}, err
	}

	return participants, nil
}

func (db *MatchStore) ListMatchHistoryIDs(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]string, error) {
	rows, err := db.Tx.Query(ctx, `
		SELECT
			Participant.match_id
		FROM
			Participant
		JOIN
			Match USING (match_id)
		WHERE
			Participant.puuid = @puuid
			AND Match.date >= @start
			AND Match.date < @end
		ORDER BY
			date DESC
		LIMIT
			100;
	`,
		pgx.NamedArgs{
			"puuid": puuid,
			"start": start,
			"end":   end,
		},
	)
	if err != nil {
		return nil, err
	}

	collect := func(row pgx.CollectableRow) (s string, err error) {
		err = row.Scan(&s)
		return s, err
	}

	ids, err := pgx.CollectRows(rows, collect)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (db * MatchStore) GetNewMatchIDs(ctx context.Context, ids []string) (newIDs []string, err error) {
	rows, err := db.Tx.Query(ctx, `
		SELECT
			match_id
		FROM
			Match
		WHERE
			match_id = any(@ids)
		ORDER BY
			date DESC;
	`,
		pgx.NamedArgs{
			"ids": ids,
		},
	)

	if err != nil {
		return nil, err
	}

	collect := func(row pgx.CollectableRow) (m string, err error) {
		err = row.Scan(&m)
		return m, err
	}

	oldIDs, err := pgx.CollectRows(rows, collect)
	if err != nil {
		return nil, err
	}

	newIDs = []string{}
	for _, id := range ids {
		found := false
		for _, oldID := range oldIDs {
			if id == oldID {
				found = true
			}
		}
		if !found {
			newIDs = append(newIDs, id)
		}
	}

	return newIDs, nil
}
