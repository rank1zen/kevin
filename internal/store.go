package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal/riot"
)

type Store struct {
	conn *pgxpool.Pool
}

func NewStore(conn *pgxpool.Pool) *Store {
	if conn == nil {
		panic("nil connection given to store")
	}
	return &Store{conn}
}

// GetItemEvents returns all item events in a match, in time order.
func (s *Store) GetItemEvents(ctx context.Context, matchID MatchID) ([]ItemEvent, error) {
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

	collect := func(row pgx.CollectableRow) (ItemEvent, error) {
		var event ItemEvent
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

func (s *Store) RecordItemEvents(ctx context.Context, events []ItemEvent) error {
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

// GetNewMatchIDs returns ids that are not in store.
func (s *Store) GetNewMatchIDs(ctx context.Context, ids []MatchID) (newIDs []MatchID, err error) {
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

	collect := func(row pgx.CollectableRow) (m MatchID, err error) {
		err = row.Scan(
			&m,
		)
		return m, err
	}

	oldIDs, err := pgx.CollectRows(rows, collect)
	if err != nil {
		return nil, err
	}

	newIDs = []MatchID{}
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

func (s *Store) ListMatchIDs(ctx context.Context, puuid string) ([]MatchID, error) {
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

	collect := func(row pgx.CollectableRow) (m MatchID, err error) {
		err = row.Scan(
			&m,
		)
		return m, err
	}

	return pgx.CollectRows(rows, collect)
}

// RecordSummoner creates new summoner and rank record, atomically.
func (s *Store) RecordSummoner(ctx context.Context, summoner Summoner, rank RankStatus) error {
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
			"platform":    summoner.Platform,
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
			"is_ranked":      rank.IsRanked,
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
				"tier":           convertTier(rank.Detail.Tier),
				"division":       convertRank(rank.Detail.Rank),
				"lp":             rank.Detail.LP,
			},
		)
		if err != nil {
			return fmt.Errorf("rank detail: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (s *Store) GetSummoner(ctx context.Context, puuid string) (Summoner, error) {
	var summoner Summoner
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
		return Summoner{}, err
	}

	return summoner, nil
}

func (s *Store) GetMatch(ctx context.Context, id MatchID) (Match, [10]Participant, error) {
	var match Match
	s.conn.QueryRow(ctx, `
		select date, duration, version, winner from Match where match_id = $1;
	`, id).Scan(&match.Date, &match.Duration, &match.Version, &match.Winner)

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

	collect := func(row pgx.CollectableRow) (m Participant, err error) {
		var runeList [11]Rune
		err = row.Scan(
			&m.Champion,
			&m.ChampionLevel,
			&m.Summoners,
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
		m.Runes = makeRunePage(runeList)
		return m, err
	}

	participants, err := pgx.CollectRows(rows, collect)
	if err != nil {
		return Match{}, [10]Participant{}, err
	}

	var p [10]Participant
	if len(participants) != 10 {
		return Match{}, [10]Participant{}, err
	} else {
		for i := range 10 {
			p[i] = participants[i]
		}
	}

	return match, p, nil
}

func (s *Store) GetChampions(ctx context.Context, puuid string) (_ []SummonerChampion, err error) {
	rows, _ := s.conn.Query(ctx, `
		select
			champion_id,
			round(avg(kills)),
			round(avg(deaths)),
			round(avg(assists)),
			avg(kill_participation),
			round(avg(creep_score)),
			avg(creep_score_per_minte),
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
		group by champion_id
		order by champion_id asc;
	`, puuid)

	collect := func(row pgx.CollectableRow) (m SummonerChampion, err error) {
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

func (s *Store) GetMatches(ctx context.Context, puuid string, page int) ([]SummonerMatch, error) {
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

	collect := func(row pgx.CollectableRow) (m SummonerMatch, err error) {
		var runeList [11]Rune
		var winner Team
		err = row.Scan(
			&m.MatchID,
			&m.Date,
			&m.Duration,
			&winner,
			&m.Team,
			&m.Champion,
			&m.ChampionLevel,
			&m.Summoners,
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
		m.Runes = makeRunePage(runeList)
		if winner == m.Team {
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

func (s *Store) GetRank(ctx context.Context, puuid string) (*RankDetail, error) {
	var m RankDetail
	var tier, rank string
	if err := s.conn.QueryRow(ctx, `
		SELECT
			detail.tier,
			detail.division,
			detail.lp,
			detail.wins,
			detail.losses
		FROM
			RankStatus status
		JOIN
			RankDetail detail USING (rank_status_id)
		WHERE
			puuid = $1 and is_current = true;
	`,
		puuid,
	).
		Scan(
			&tier,
			&rank,
			&m.LP,
			&m.Wins,
			&m.Losses,
		); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	m.Tier = convertTier2(tier)
	m.Rank = convertRank2(rank)

	return &m, nil
}

func (s *Store) SearchSummoner(ctx context.Context, q string) (_ []SearchResult, err error) {
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

	collect := func(row pgx.CollectableRow) (m SearchResult, err error) {
		err = row.Scan(&m.Puuid, &m.Name, &m.Tagline)
		return m, err
	}

	results, err := pgx.CollectRows(rows, collect)
	return results, err
}

func (s *Store) RecordMatch(ctx context.Context, match Match, participants [10]Participant) error {
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
			"winner":   match.Winner,
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
			participant.Puuid,
			participant.Team,
			participant.Champion,
			participant.ChampionLevel,
			participant.Summoners,
			makeRuneList(participant.Runes),
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

func makeRuneList(runes RunePage) [11]int {
	ids := [11]Rune{
		runes.PrimaryTree,
		runes.PrimaryKeystone,
		runes.PrimaryA,
		runes.PrimaryB,
		runes.PrimaryC,
		runes.SecondaryTree,
		runes.SecondaryA,
		runes.SecondaryB,
		runes.MiniOffense,
		runes.MiniFlex,
		runes.MiniDefense,
	}

	res := [11]int{}
	for i := range 11 {
		res[i] = int(ids[i])
	}
	return res
}

func convertTier(tier riot.Tier) string {
	switch tier {
	default:
		panic("bro.")
	case riot.TierIron:
		return "Iron"
	case riot.TierBronze:
		return "Bronze"
	case riot.TierSilver:
		return "Silver"
	case riot.TierGold:
		return "Gold"
	case riot.TierPlatinum:
		return "Platinum"
	case riot.TierEmerald:
		return "Emerald"
	case riot.TierDiamond:
		return "Diamond"
	case riot.TierMaster:
		return "Master"
	case riot.TierGrandmaster:
		return "Grandmaster"
	case riot.TierChallenger:
		return "Challenger"
	}
}

func convertRank(rank riot.Rank) string {
	switch rank {
	default:
		panic("bro.")
	case riot.Rank1:
		return "I"
	case riot.Rank2:
		return "II"
	case riot.Rank3:
		return "III"
	case riot.Rank4:
		return "IV"
	}
}

func convertTier2(tier string) riot.Tier {
	switch tier {
	default:
		panic("bro.")
	case "Iron":
		return riot.TierIron
	case "Bronze":
		return riot.TierBronze
	case "Silver":
		return riot.TierSilver
	case "Gold":
		return riot.TierGold
	case "Platinum":
		return riot.TierPlatinum
	case "Emerald":
		return riot.TierEmerald
	case "Diamond":
		return riot.TierDiamond
	case "Master":
		return riot.TierMaster
	case "Grandmaster":
		return riot.TierGrandmaster
	case "Challenger":
		return riot.TierChallenger
	}
}

func convertRank2(rank string) riot.Rank {
	switch rank {
	default:
		panic("bro.")
	case "I":
		return riot.Rank1
	case "II":
		return riot.Rank2
	case "III":
		return riot.Rank3
	case "IV":
		return riot.Rank4
	}
}
