package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/pgxutil"
)

func createMatchInfo(ctx context.Context, conn pgxutil.Query, m internal.RiotMatch) (internal.RiotMatch, error) {
	row := pgx.NamedArgs{
		"match_id":     m.ID,
		"data_version": m.DataVersion,
		"date":         m.EndTimestamp,
		"duration":     m.Duration,
		"patch":        m.Patch,
	}

	var result internal.RiotMatch

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

func createParticipant(ctx context.Context, conn pgxutil.Exec, m internal.RiotMatchParticipant) error {
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

func createTeam(ctx context.Context, conn pgxutil.Exec, m internal.RiotMatchTeam) error {
	_, err := conn.Exec(ctx, `
	INSERT INTO match_teams
		(match_id, team_id, win)
	VALUES
		($1, $2, $3)
	`, m.MatchID, m.ID, m.Win)

	return err
}

func (db *DB) CreateMatch(ctx context.Context, match internal.RiotMatch) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	// note the other two tables depend on this one
	createMatchInfo(ctx, tx, match)
	if err != nil {
		return err
	}

	for _, participant := range match.GetParticipants() {
		err := createParticipant(ctx, tx, participant)
		if err != nil {
			return err
		}
	}

	for _, team := range match.GetTeams() {
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

// GetMatchHistory returns a paging object for a summoner's match history.
func (db *DB) GetMatchHistory(ctx context.Context, puuid internal.PUUID) (internal.MatchHistory, error) {
	rows, _ := db.pool.Query(ctx, `
	`, puuid)

	collectFn := func(row pgx.CollectableRow) (internal.Participant, error) {
		var p internal.Participant
		err := row.Scan(
			&p.Puuid,
			&p.ID,
			&p.MatchID,
			&p.TeamID,
			&p.StartTimestamp,
			&p.EndTimestamp,
			&p.Duration,
			&p.Patch,
			&p.Position,
			&p.Win,
			&p.BannedChampion,
			&p.Champion,
			&p.ChampionLevel,
			&p.Summs,
			&p.Items,
			&p.Runes,
			&p.Kills,
			&p.Deaths,
			&p.Assists,
			&p.KillParticipation,
			&p.CreepScore,
			&p.CsPerMinute,
			&p.Gold,
			&p.GoldPercentage,
			&p.GoldDelta,
			&p.Damage,
			&p.DamagePercentage,
			&p.DamageDelta,
			&p.VisionScore,
			&p.PinkWards,
		)

		return p, err
	}

	matches, err := pgx.CollectRows(rows, collectFn)
	if err != nil {
		return internal.MatchHistory{}, err
	}

	return internal.MatchHistory{
		List: matches,
	}, nil
}

// GetMatchHistory returns the 10 participants in a match.
func (db *DB) GetMatch(ctx context.Context, matchID internal.MatchID) ([10]internal.Participant, error) {
	rows, _ := db.pool.Query(ctx, `
	`, matchID)

	collectFn := func(row pgx.CollectableRow) (internal.Participant, error) {
		var p internal.Participant
		err := row.Scan(
			&p.Puuid,
			&p.ID,
			&p.MatchID,
			&p.TeamID,
			&p.StartTimestamp,
			&p.EndTimestamp,
			&p.Duration,
			&p.Patch,
			&p.Position,
			&p.Win,
			&p.BannedChampion,
			&p.Champion,
			&p.ChampionLevel,
			&p.Summs,
			&p.Items,
			&p.Runes,
			&p.Kills,
			&p.Deaths,
			&p.Assists,
			&p.KillParticipation,
			&p.CreepScore,
			&p.CsPerMinute,
			&p.Gold,
			&p.GoldPercentage,
			&p.GoldDelta,
			&p.Damage,
			&p.DamagePercentage,
			&p.DamageDelta,
			&p.VisionScore,
			&p.PinkWards,
		)

		return p, err
	}

	matches, err := pgx.CollectRows(rows, collectFn)
	if err != nil {
		return [10]internal.Participant{}, err
	}

	// make sure there are ten matches

	return [10]internal.Participant{matches[0]}, nil // FIXME
}
