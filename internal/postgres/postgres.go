package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal"
)

// Store manages connections with a postgres database.
type Store struct{
	conn *pgxpool.Pool
}

// TODO: postgres package should be responsible for creating connections to postgres
func NewStore(conn *pgxpool.Pool) *Store {
	if conn == nil {
		panic("nil connection given to store")
	}
	return &Store{conn}
}

func (s *Store) GetPUUID(ctx context.Context, name, tag string) (puuid string, err error) {
	err = s.conn.QueryRow(ctx, `
		SELECT puuid FROM summoner WHERE name = $1 AND tagline = $2;
	`, name, tag).Scan(&puuid)

	return puuid, err
}

func (s *Store) GetZMatches(ctx context.Context, puuid string, start, end time.Time) ([]internal.SummonerMatch, error) {
	rows, _ := s.conn.Query(ctx, `
		SELECT
			m.match_id,
			m.date,
			m.duration,
			m.winner,
			p.team,
			p.champion,
			p.champion_level,
			p.summoners,
			p.runes,
			p.items,
			p.kills,
			p.deaths,
			p.assists,
			p.kill_participation,
			p.creep_score,
			p.creep_score_per_minute,
			p.gold_earned,
			p.gold_delta_enemy,
			p.gold_percentage_team,
			p.damage_dealt,
			p.damage_taken,
			p.damage_delta_enemy,
			p.damage_percentage_team,
			p.vision_score,
			p.pink_wards_bought
		FROM
			Participant as p
		JOIN
			Match as m USING (match_id)
		WHERE
			puuid = @puuid
		AND
			m.date >= @start
		AND
			m.date <= @end
		ORDER BY
			m.date desc
	`,
		pgx.NamedArgs{
			"puuid": puuid,
			"start": start,
			"end": end,
		},
	)

	collect := func(row pgx.CollectableRow) (m internal.SummonerMatch, err error) {
		var runeList [11]int
		var winner int
		err = row.Scan(
			&m.MatchID,
			&m.Date,
			&m.Duration,
			&winner,
			&m.TeamID,
			&m.ChampionID,
			&m.ChampionLevel,
			&m.SummonerIDs,
			&runeList,
			&m.Items,
			&m.Kills,
			&m.Deaths,
			&m.Assists,
			&m.KillParticipation,
			&m.CreepScore,
			&m.CreepScorePerMinute,
			&m.GoldEarned,
			&m.GoldDeltaEnemy,
			&m.GoldPercentageTeam,
			&m.DamageDealt,
			&m.DamageTaken,
			&m.DamageDeltaEnemy,
			&m.DamagePercentageTeam,
			&m.VisionScore,
			&m.PinkWardsBought,
		)
		m.Runes = internal.NewRunePage(internal.WithIntList(runeList))
		if winner == m.TeamID {
			m.Win = true
		} else {
			m.Win = false
		}
		return m, err
	}

	return pgx.CollectRows(rows, collect)
}

func (s *Store) RecordTimeline(ctx context.Context, id string, items []internal.ItemEvent, skills []internal.SkillEvent) error {
	var batch pgx.Batch

	for _, event := range items {
		batch.Queue(`
			INSERT INTO ItemEvent (
				match_id,
				puuid,
				in_game_timestamp,
				item_id,
				type
			)
			VALUES (
				$1,
				$2,
				$3,
				$4,
				$5
			)
		`,
			event.MatchID,
			event.PUUID,
			event.InGameTimestamp,
			event.ItemID,
			event.Type,
		)
	}

	for _, event := range skills {
		batch.Queue(`
			INSERT INTO SkillEvent (
				match_id,
				puuid,
				in_game_timestamp,
				spell_slot
			)
			VALUES (
				$1,
				$2,
				$3,
				$4
			)
		`,
			event.MatchID,
			event.PUUID,
			event.InGameTimestamp,
			event.SpellSlot,
		)
	}

	return s.conn.SendBatch(ctx, &batch).Close()
}

// GetItemEvents returns all item events in a match, in time order.
func (s *Store) GetItemEvents(ctx context.Context, matchID string) ([]internal.ItemEvent, error) {
	rows, _ := s.conn.Query(ctx, `
		SELECT
			match_id,
			puuid,
			in_game_timestamp,
			item_id,
			type
		FROM
			ItemEvent
		WHERE
			match_id = $1
		ORDER BY
			in_game_timestamp ASC;
	`,
		matchID)

	collect := func(row pgx.CollectableRow) (internal.ItemEvent, error) {
		var event internal.ItemEvent
		err := row.Scan(
			&event.MatchID,
			&event.PUUID,
			&event.InGameTimestamp,
			&event.ItemID,
			&event.Type,
		)
		return event, err
	}

	return pgx.CollectRows(rows, collect)
}

func (s *Store) recordItemEvents(ctx context.Context, events []internal.ItemEvent) error {
	batch := pgx.Batch{}

	for _, event := range events {
		batch.Queue(`
			INSERT INTO ItemEvent (
				match_id,
				puuid,
				in_game_timestamp,
				item_id,
				type
			)
			VALUES (
				$1,
				$2,
				$3,
				$4,
				$5
			)
		`,
			event.MatchID,
			event.PUUID,
			event.InGameTimestamp,
			event.ItemID,
			event.Type,
		)
	}

	return s.conn.SendBatch(ctx, &batch).Close()
}

func (s *Store) recordSkillEvents(ctx context.Context, events []internal.SkillEvent) error {
	batch := pgx.Batch{}

	for _, event := range events {
		batch.Queue(`
			INSERT INTO SkillEvent (
				match_id,
				puuid,
				in_game_timestamp,
				spell_slot
			)
			VALUES (
				$1,
				$2,
				$3,
				$4
			)
		`,
			event.MatchID,
			event.PUUID,
			event.InGameTimestamp,
			event.SpellSlot,
		)
	}

	return s.conn.SendBatch(ctx, &batch).Close()
}

// GetNewMatchIDs returns ids that are not in store.
func (s *Store) GetNewMatchIDs(ctx context.Context, ids []string) (newIDs []string, err error) {
	rows, err := s.conn.Query(ctx, `
		SELECT
			match_id
		FROM
			Match
		WHERE
			match_id = any($1)
		ORDER BY
			date DESC;
	`,
		ids,
	)
	if err != nil {
		return nil, err
	}

	collect := func(row pgx.CollectableRow) (m string, err error) {
		err = row.Scan(
			&m,
		)
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

func (s *Store) ListMatchIDs(ctx context.Context, puuid string) ([]string, error) {
	rows, err := s.conn.Query(ctx, `
		SELECT
			match_id
		FROM
			Participant
		WHERE
			puuid = $1
	`,
		puuid,
	)
	if err != nil {
		return nil, err
	}

	collect := func(row pgx.CollectableRow) (m string, err error) {
		err = row.Scan(
			&m,
		)
		return m, err
	}

	return pgx.CollectRows(rows, collect)
}

func (s *Store) RecordSummoner(ctx context.Context, summoner internal.Summoner, rank internal.RankStatus) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO Summoner (
			puuid,
			platform,
			name,
			tagline,
			summoner_id
		)
		VALUES (
			@puuid,
			@platform,
			@name,
			@tagline,
			@summoner_id
		)
		ON CONFLICT (puuid)
		DO UPDATE SET
			name = @name,
			tagline = @tagline,
			platform = @platform,
			summoner_id = @summoner_id;
	`,
		pgx.NamedArgs{
			"puuid":       summoner.PUUID,
			"platform":    "NA1",
			"name":        summoner.Name,
			"tagline":     summoner.Tagline,
			"summoner_id": summoner.SummonerID,
		},
	)
	if err != nil {
		return fmt.Errorf("summoner: %w", err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE RankStatus SET
			end_date = @effective_date,
			is_current = false
		WHERE
			is_current = true;
	`,
		pgx.NamedArgs{
			"effective_date": rank.EffectiveDate,
		},
	)
	if err != nil {
		return fmt.Errorf("updating rank status: %w", err)
	}

	var recordID int

	var isRanked bool
	if rank.Detail != nil {
		isRanked = true
	} else {
		isRanked = false
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO RankStatus (
			puuid,
			effective_date,
			end_date,
			is_current,
			is_ranked
		)
		VALUES (
			@puuid,
			@effective_date,
			'infinity',
			true,
			@is_ranked
		)
		RETURNING
			rank_status_id;
	`,
		pgx.NamedArgs{
			"puuid":          rank.PUUID,
			"effective_date": rank.EffectiveDate,
			"is_ranked":      isRanked,
		},
	).Scan(&recordID)
	if err != nil {
		return fmt.Errorf("rank status: %w", err)
	}

	if rank.Detail != nil {
		_, err := tx.Exec(ctx, `
			INSERT INTO RankDetail (
				rank_status_id,
				wins,
				losses,
				tier,
				division,
				lp
			)
			VALUES (
				@rank_status_id,
				@wins,
				@losses,
				@tier,
				@division,
				@lp
			)
		`,
			pgx.NamedArgs{
				"rank_status_id": recordID,
				"wins":           rank.Detail.Wins,
				"losses":         rank.Detail.Losses,
				"tier":           convertRiotTierToString(rank.Detail.Tier),
				"division":       rank.Detail.Division,
				"lp":             rank.Detail.LP,
			},
		)
		if err != nil {
			return fmt.Errorf("rank detail: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (s *Store) GetSummoner(ctx context.Context, puuid string) (internal.Summoner, error) {
	var summoner internal.Summoner
	if err := s.conn.QueryRow(ctx, `
		SELECT
			puuid,
			name,
			tagline,
			platform,
			summoner_id
		FROM
			Summoner
		WHERE
			puuid = $1;
	`, puuid).Scan(
		&summoner.PUUID,
		&summoner.Name,
		&summoner.Tagline,
		&summoner.Platform,
		&summoner.SummonerID,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return internal.Summoner{}, internal.ErrSummonerNotFound
		}
		return internal.Summoner{}, err
	}

	return summoner, nil
}

func (s *Store) GetMatch(ctx context.Context, id string) (internal.Match, [10]internal.Participant, error) {
	var match internal.Match
	s.conn.QueryRow(ctx, `
		select date, duration, version, winner from Match where match_id = $1;
	`, id).Scan(&match.Date, &match.Duration, &match.Version, &match.WinnerID)

	rows, _ := s.conn.Query(ctx, `
		select
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
		from Participant where match_id = $1;
	`, id)

	collect := func(row pgx.CollectableRow) (m internal.Participant, err error) {
		var runeList [11]int
		err = row.Scan(
			&m.ChampionID,
			&m.ChampionLevel,
			&m.SummonerIDs,
			&runeList,
			&m.Items,
			&m.Kills,
			&m.Deaths,
			&m.Assists,
			&m.KillParticipation,
			&m.CreepScore,
			&m.CreepScorePerMinute,
			&m.GoldEarned,
			&m.GoldDeltaEnemy,
			&m.GoldPercentageTeam,
			&m.DamageDealt,
			&m.DamageTaken,
			&m.DamageDeltaEnemy,
			&m.DamagePercentageTeam,
			&m.VisionScore,
			&m.PinkWardsBought,
		)
		m.Runes = internal.NewRunePage(internal.WithIntList(runeList))
		return m, err
	}

	participants, err := pgx.CollectRows(rows, collect)
	if err != nil {
		return internal.Match{}, [10]internal.Participant{}, err
	}

	var p [10]internal.Participant
	if len(participants) != 10 {
		return internal.Match{}, [10]internal.Participant{}, err
	} else {
		for i := range 10 {
			p[i] = participants[i]
		}
	}

	return match, p, nil
}

func (s *Store) GetChampions(ctx context.Context, puuid string) (_ []internal.SummonerChampion, err error) {
	rows, _ := s.conn.Query(ctx, `
		select
			champion,
			round(avg(kills)),
			round(avg(deaths)),
			round(avg(assists)),
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
		from Participant
		where puuid = $1
		group by champion
		order by champion asc;
	`, puuid)

	collect := func(row pgx.CollectableRow) (m internal.SummonerChampion, err error) {
		err = row.Scan(
			&m.Champion,
			&m.Kills,
			&m.Deaths,
			&m.Assists,
			&m.KillParticipation,
			&m.CreepScore,
			&m.CreepScorePerMinute,
			&m.GoldEarned,
			&m.GoldDeltaEnemy,
			&m.GoldPercentageTeam,
			&m.DamageDealt,
			&m.DamageDeltaEnemy,
			&m.DamagePercentageTeam,
			&m.DamageTaken,
			&m.VisionScore,
			&m.PinkWardsBought,
		)
		return m, err
	}

	champions, err := pgx.CollectRows(rows, collect)
	return champions, err
}

func (s *Store) GetMatches(ctx context.Context, puuid string, page int) ([]internal.SummonerMatch, error) {
	rows, _ := s.conn.Query(ctx, `
		SELECT
			m.match_id,
			m.date,
			m.duration,
			m.winner,
			p.team,
			p.champion,
			p.champion_level,
			p.summoners,
			p.runes,
			p.items,
			p.kills,
			p.deaths,
			p.assists,
			p.kill_participation,
			p.creep_score,
			p.creep_score_per_minute,
			p.gold_earned,
			p.gold_delta_enemy,
			p.gold_percentage_team,
			p.damage_dealt,
			p.damage_taken,
			p.damage_delta_enemy,
			p.damage_percentage_team,
			p.vision_score,
			p.pink_wards_bought
		FROM
			Participant as p
		JOIN
			Match as m using (match_id)
		WHERE
			puuid = $1
		ORDER BY
			m.date desc
		OFFSET
			$2
		LIMIT
			$3;
	`, puuid, page*10, 10)

	collect := func(row pgx.CollectableRow) (m internal.SummonerMatch, err error) {
		var runeList [11]int
		var winner int
		err = row.Scan(
			&m.MatchID,
			&m.Date,
			&m.Duration,
			&winner,
			&m.TeamID,
			&m.ChampionID,
			&m.ChampionLevel,
			&m.SummonerIDs,
			&runeList,
			&m.Items,
			&m.Kills,
			&m.Deaths,
			&m.Assists,
			&m.KillParticipation,
			&m.CreepScore,
			&m.CreepScorePerMinute,
			&m.GoldEarned,
			&m.GoldDeltaEnemy,
			&m.GoldPercentageTeam,
			&m.DamageDealt,
			&m.DamageTaken,
			&m.DamageDeltaEnemy,
			&m.DamagePercentageTeam,
			&m.VisionScore,
			&m.PinkWardsBought,
		)
		m.Runes = internal.NewRunePage(internal.WithIntList(runeList))
		if winner == m.TeamID {
			m.Win = true
		} else {
			m.Win = false
		}
		return m, err
	}

	matches, err := pgx.CollectRows(rows, collect)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func (s *Store) GetRank(ctx context.Context, puuid string, ts time.Time, recent bool) (m internal.RankRecord, err error) {
	args := pgx.NamedArgs{
		"puuid": puuid,
		"ts": ts,
	}

	genQuery := func(s string) string {
		return `
			SELECT
				rank_status_id,
				puuid,
				effective_date,
				end_date,
				is_current,
				is_ranked
			FROM
				RankStatus status
			WHERE
				puuid = @puuid
		` + s + `
			LIMIT
				1;
		`
	}

	var q string
	if recent {
		q = genQuery(`
				AND status.effective_date <= @ts
			ORDER BY
				status.effective_date DESC
		`)
	} else {
		q = genQuery(`
				AND status.effective_date > @ts
			ORDER BY
				status.effective_date ASC
		`)
	}

	var (
		statusID int
		timestamp pgtype.Timestamp
		isRanked bool
	)

	if err = s.conn.QueryRow(ctx, q, args).Scan(
		&statusID,
		&m.PUUID,
		&m.EffectiveDate,
		&timestamp,
		&m.IsCurrent,
		&isRanked,
	); err != nil {
		return m, err
	}

	if timestamp.InfinityModifier == pgtype.Finite {
		m.EndDate = &timestamp.Time
	}

	if isRanked {
		details, err := s.getRankDetail(ctx, statusID)
		if err != nil {
			return m, err
		}

		m.Detail = &details
		return m, nil
	}

	return m, nil
}

func (s *Store) SearchSummoner(ctx context.Context, q string) (_ []internal.SearchResult, err error) {
	rows, _ := s.conn.Query(ctx, `
		WITH rankings AS (
			SELECT
				puuid,
				name,
				tagline,
				to_tsvector(name) as txt,
				websearch_to_tsquery($1) as query
			FROM
				Summoner
		)
		SELECT
			puuid,
			name,
			tagline
		FROM
			rankings
		WHERE
			txt @@ query
		ORDER BY
			ts_rank(txt, query)
		LIMIT 10;
	`, q)

	collect := func(row pgx.CollectableRow) (m internal.SearchResult, err error) {
		err = row.Scan(&m.Puuid, &m.Name, &m.Tagline)
		return m, err
	}

	results, err := pgx.CollectRows(rows, collect)
	return results, err
}

func (s *Store) RecordMatch(ctx context.Context, match internal.Match, participants [10]internal.Participant) error {
	batch := pgx.Batch{}

	batch.Queue(`
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

	for _, participant := range participants {
		batch.Queue(`
			INSERT INTO Participant (
				match_id,
				puuid,
				team,
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
				$23
			);
		`,
			participant.MatchID,
			participant.PUUID,
			participant.TeamID,
			participant.ChampionID,
			participant.ChampionLevel,
			participant.SummonerIDs,
			convertRunePageToList(participant.Runes),
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

	return s.conn.SendBatch(ctx, &batch).Close()
}

func (s *Store) getRankDetail(ctx context.Context, statusID int) (m internal.RankDetail, err error) {
	var (
		tier string
		rank string
	)

	err = s.conn.QueryRow(ctx, `
		SELECT
			tier,
			division,
			lp,
			wins,
			losses
		FROM
			RankDetail
		WHERE
			rank_status_id = $1;
	`,
		statusID,
	).Scan(
		&tier,
		&rank,
		&m.LP,
		&m.Wins,
		&m.Losses,
	)

	m.Tier = convertStringToRiotTier(tier)
	m.Division = convertStringToRiotRank(rank)

	return m, err
}
