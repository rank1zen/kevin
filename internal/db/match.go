package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/pgxutil"
)

func createMatch(ctx context.Context, conn pgxutil.Query, m internal.Match) (internal.Match, error) {
	row := pgx.NamedArgs{
		"match_id":     m.ID,
		"data_version": m.DataVersion,
		"date":         m.EndTimestamp,
		"duration":     m.Duration,
		"patch":        m.Patch,
	}

	var result internal.Match

	err := conn.QueryRow(ctx, `
	INSERT INTO matches (
		match_id,
		data_version,
		date,
		duration,
		patch
	)
	VALUES (
		@match_id,
		@data_version,
		@date,
		@duration,
		@patch
	)
	RETURNING
		match_id,
		data_version,
		date,
		duration,
		patch;
	`, row).Scan(
		&result.ID,
		&result.DataVersion,
		&result.EndTimestamp,
		&result.Duration,
		&result.Patch,
	)

	return result, err
}

func createParticipant(ctx context.Context, conn pgxutil.Exec, m internal.MatchParticipant) error {
	row := pgx.NamedArgs{
		"match_id":                           m.Match,
		"participant_id":                     m.ID,
		"puuid":                              m.Puuid,
		"team_id":                            m.Team,
		"kills":                              m.Kills,
		"assists":                            m.Assists,
		"deaths":                             m.Deaths,
		"creep_score":                        m.TotalMinionsKilled + m.NeutralMinionsKilled,
		"vision_score":                       m.VisionScore,
		"gold_earned":                        m.GoldEarned,
		"gold_spent":                         m.GoldSpent,
		"player_position":                    m.Role,
		"champion_level":                     m.ChampionLevel,
		"champion_id":                        m.ChampionID,
		"champion_name":                      m.ChampionName,
		"items":                              m.Items,
		"summoners":                          m.Summoners,
		"runes":                              m.Runes.ToList(),
		"physical_damage_dealt":              m.PhysicalDamageDealt,
		"physical_damage_dealt_to_champions": m.PhysicalDamageDealtToChampions,
		"physical_damage_taken":              m.PhysicalDamageTaken,
		"magic_damage_dealt":                 m.MagicDamageDealt,
		"magic_damage_dealt_to_champions":    m.MagicDamageDealtToChampions,
		"magic_damage_taken":                 m.MagicDamageTaken,
		"true_damage_dealt":                  m.TrueDamageDealt,
		"true_damage_dealt_to_champions":     m.TrueDamageDealtToChampions,
		"true_damage_taken":                  m.TrueDamageTaken,
		"total_damage_dealt":                 m.TotalDamageDealt,
		"total_damage_dealt_to_champions":    m.TotalDamageDealtToChampions,
		"total_damage_taken":                 m.TotalDamageTaken,
	}

	_, err := conn.Exec(ctx, `
	INSERT INTO match_participants (
		match_id,
		participant_id,
		puuid,
		team_id,
		kills,
		assists,
		deaths,
		creep_score,
		vision_score,
		gold_earned,
		gold_spent,
		player_position,
		champion_level,
		champion_id,
		champion_name,
		items,
		summoners,
		runes,
		physical_damage_dealt,
		physical_damage_dealt_to_champions,
		physical_damage_taken,
		magic_damage_dealt,
		magic_damage_dealt_to_champions,
		magic_damage_taken,
		true_damage_dealt,
		true_damage_dealt_to_champions,
		true_damage_taken,
		total_damage_dealt,
		total_damage_dealt_to_champions,
		total_damage_taken
	)
	VALUES (
		@match_id,
		@participant_id,
		@puuid,
		@team_id,
		@kills,
		@assists,
		@deaths,
		@creep_score,
		@vision_score,
		@gold_earned,
		@gold_spent,
		@player_position,
		@champion_level,
		@champion_id,
		@champion_name,
		@items,
		@summoners,
		@runes,
		@physical_damage_dealt,
		@physical_damage_dealt_to_champions,
		@physical_damage_taken,
		@magic_damage_dealt,
		@magic_damage_dealt_to_champions,
		@magic_damage_taken,
		@true_damage_dealt,
		@true_damage_dealt_to_champions,
		@true_damage_taken,
		@total_damage_dealt,
		@total_damage_dealt_to_champions,
		@total_damage_taken
	);
	`, row)

	return err
}

func createTeam(ctx context.Context, conn pgxutil.Exec, m internal.MatchTeam) error {
	_, err := conn.Exec(ctx, `
	INSERT INTO match_teams
		(match_id, team_id, win)
	VALUES
		($1, $2, $3)
	`, m.MatchID, m.ID, m.Win)

	return err
}

func MatchInsert(ctx context.Context, conn pgxutil.Conn, m internal.Match) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	// note the other two tables depend on this one
	createMatch(ctx, conn, m)
	if err != nil {
		return err
	}

	for _, participant := range m.GetParticipants() {
		err := createParticipant(ctx, tx, participant)
		if err != nil {
			return err
		}
	}

	for _, team := range m.GetTeams() {
		err := createTeam(ctx, tx, team)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
