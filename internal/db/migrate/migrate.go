package migrate

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func Migrate(ctx context.Context, conn *pgx.Conn) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, f := range migrations {
		err := f(tx)
		if err != nil {
			return fmt.Errorf("[migration %d]: %w", i, err)
		}
	}

	return tx.Commit(ctx)
}

var schemaVersion = len(migrations)

var migrations = []func(tx pgx.Tx) error{
	func(tx pgx.Tx) (err error) {
		sql := `
		CREATE TABLE schema_version (
			version text not null
		);

		CREATE DOMAIN riot_puuid       AS CHAR(78);
		CREATE DOMAIN riot_summoner_id AS VARCHAR(63);
		CREATE DOMAIN riot_account_id  AS VARCHAR(56);
		CREATE DOMAIN riot_match_id    AS VARCHAR(60);

		CREATE TABLE profiles (
			puuid        riot_puuid  primary key,
			name         varchar(32) not null,
			tagline      varchar(10) not null,
			last_updated timestamptz not null
		);

		CREATE TABLE summoner_records (
			record_id uuid default gen_random_uuid() primary key,
			valid_from      timestamp        not null,
			valid_to        timestamp        not null,
			entered_at      timestamp        not null,
			account_id      riot_account_id  not null,
			summoner_id     riot_summoner_id not null,
			puuid           riot_puuid       not null,
			revision_date   timestamptz      not null,
			summoner_level  bigint           not null,
			profile_icon_id int              not null
		);

		CREATE TABLE league_records (
			record_id uuid default gen_random_uuid() primary key,
			summoner_id   riot_summoner_id not null,
			valid_from    timestamp        not null,
			valid_to      timestamp        not null,
			entered_at    timestamp        not null,
			is_ranked     boolean          not null,
			league_id     varchar(128),
			tier          varchar(16),
			division      varchar(8),
			league_points int,
			wins          int,
			losses        int
		);

		CREATE TABLE matchlist (
			puuid    riot_puuid    not null,
			match_id riot_match_id not null,
			group    uuid          not null,
			unique (match_id, puuid)
		);

		CREATE TABLE matches (
			match_id     riot_match_id primary key,
			data_version text          not null,
			date         timestamptz   not null,
			duration     interval      not null,
			patch        text          not null
		);

		CREATE TABLE match_participants (
			match_id riot_match_id not null,
			participant_id int not null,
			puuid riot_puuid  not null,
			FOREIGN KEY (match_id)
				REFERENCES matches (match_id)
				ON DELETE CASCADE,
			UNIQUE (match_id, participant_id),
			UNIQUE (match_id, puuid),
			team_id int not null,

			name            varchar(50) not null,
			player_position varchar(10) not null,
			champion_level  int         not null,
			champion_id     int         not null,
			champion_name   varchar(30) not null,
			kills           int         not null,
			deaths          int         not null,
			assists         int         not null,
			creep_score     int         not null,
			vision_score    int         not null,
			gold_earned     int         not null,
			gold_spent      int         not null,
			items           int[7]      not null,
			summoners       int[2]      not null,
			runes           int[11]     not null,

			physical_damage_dealt              int not null,
			physical_damage_dealt_to_champions int not null,
			physical_damage_taken              int not null,
			magic_damage_dealt                 int not null,
			magic_damage_dealt_to_champions    int not null,
			magic_damage_taken                 int not null,
			true_damage_dealt                  int not null,
			true_damage_dealt_to_champions     int not null,
			true_damage_taken                  int not null,
			total_damage_dealt                 int not null,
			total_damage_dealt_to_champions    int not null,
			total_damage_taken                 int not null
		);

		CREATE VIEW match_participant_deltas AS
		SELECT
			me.match_id,
			me.team_id,
			me.id,
			me.puuid,
			other.puuid,
			me.kills - other.kills,
			me.deaths - other.deaths,
			me.assists - other.assists,
			me.creep_score - other.creep_score,
			me.vision_score - other.vision_score,
			me.gold_earned - other.gold_earned,
			me.total_damage_dealt_to_champions - other.total_damage_dealt_to_champions
		FROM
			match_participants me
		JOIN
			match_participants other USING (match_id);

		CREATE VIEW match_team_stats AS
		SELECT
			match_id,
			team_id,
			sum(kills),
			sum(deaths),
			sum(assists),
			sum(creep_score),
			sum(vision_score),
			sum(gold_earned),
			sum(total_damage_dealt_to_champions)
		FROM
			match_participants
		GROUP BY
			match_id, team_id;

		CREATE VIEW match_participant_stats AS
		WITH team_totals AS (
			SELECT
				match_id,
				team_id,

				sum(kills)                           AS kills,
				sum(deaths)                          AS deaths,
				sum(assists)                         AS assists,
				sum(creep_score)                     AS creep_score,
				sum(vision_score)                    AS vision_score,
				sum(gold_earned)                     AS gold_earned,
				sum(total_damage_dealt_to_champions) AS total_damage_dealt_to_champions
			FROM
				match_participants
			GROUP BY
				match_id, team_id
		)
		WITH counterpart AS (
			SELECT
				me.match_id AS match_id,
				me.puuid    AS puuid,

				me.kills - other.kills                                                     AS kills,
				me.deaths - other.deaths                                                   AS deaths,
				me.assists - other.assists                                                 AS assists,
				me.creep_score - other.creep_score                                         AS creep_score,
				me.vision_score - other.vision_score                                       AS vision_score,
				me.gold_earned - other.gold_earned                                         AS gold_earned,
				me.total_damage_dealt_to_champions - other.total_damage_dealt_to_champions AS total_damage_dealt_to_champions
			FROM
				match_participants AS me
			JOIN
				match_participants AS other USING(match_id, participant_id)
		)
		SELECT
			kills,
			deaths,
			assists,
			(kills + assists) / team_totals.kills                                      AS kill_participation,
			creep_score,
			creep_score / (extract(epoch from game_duration) / 60)                     AS cs_per_minute,
			gold_earned,
			gold_earned / team_totals.gold_earned                                      AS gold_percentage,
			counterpart.gold_earned                                                    AS gold_delta,
			total_damage_dealt_to_champions,
			total_damage_dealt_to_champions / (extract(epoch from game_duration) / 60) AS damage_per_minute,
			total_damage_dealt_to_champions / team_totals.damage                       AS damage_percentage,
			counterpart.damage                                                         AS damage_delta,
			vision_score
		FROM
			match_participants
		JOIN
			team_totals USING (match_id, team_id)
		JOIN
			counterpart USING (match_id, puuid);

		CREATE TABLE match_teams (
			match_id riot_match_id not null,
			FOREIGN KEY (match_id)
				REFERENCES matches(id)
				ON DELETE CASCADE,
			id int not null,
			unique (match_id, id),
			win boolean not null
		);

		CREATE VIEW summoner_records_latest AS
		SELECT DISTINCT ON (puuid)
			summoner_id,
			puuid,
			summoner_level,
			profile_icon_id
		FROM
			summoner_records
		ORDER BY
			puuid,
			record_date DESC;

		CREATE VIEW league_records_latest AS
		SELECT DISTINCT ON (summoner_id)
			summoner_id,
			tier,
			division,
			league_points,
			wins,
			losses
		FROM
			league_records
		ORDER BY
			summoner_id,
			record_date DESC;
		`
		_, err = tx.Exec(context.Background(), sql)
		return err
	},
}
