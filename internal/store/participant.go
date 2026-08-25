package store

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Participant struct {
	ID              string   `db:"id"`
	MatchID         string   `db:"match_id"`
	PUUID           string   `db:"puuid"`
	TeamID          string   `db:"team_id"`
	TeamPosition    string   `db:"team_position"`
	ChampionID      string   `db:"champion_id"`
	ChampionLevel   int      `db:"champion_level"`
	SummonerIDs     []string `db:"summoner_ids"`
	RuneIDs         []string `db:"rune_ids"`
	ItemIDs         []string `db:"item_ids"`
	Kills           int      `db:"kills"`
	Deaths          int      `db:"deaths"`
	Assists         int      `db:"assists"`
	CreepScore      int      `db:"creep_score"`
	DamageDealt     int      `db:"damage_dealt"`
	DamageTaken     int      `db:"damage_taken"`
	GoldEarned      int      `db:"gold_earned"`
	VisionScore     int      `db:"vision_score"`
	PinkWardsBought int      `db:"pink_wards_bought"`
}

type CreateParticipant struct {
	MatchID         string
	PUUID           string
	TeamID          string
	TeamPosition    string
	ChampionID      string
	ChampionLevel   int
	SummonerIDs     []string
	RuneIDs         []string
	ItemIDs         []string
	Kills           int
	Deaths          int
	Assists         int
	CreepScore      int
	DamageDealt     int
	DamageTaken     int
	GoldEarned      int
	VisionScore     int
	PinkWardsBought int
}

type ParticipantStore struct {
	tx DBTX
}

func NewParticipantStore(tx DBTX) *ParticipantStore {
	return &ParticipantStore{
		tx: tx,
	}
}

func (s *ParticipantStore) GetParticipantByMatchIDAndPUUID(ctx context.Context, matchID, puuid string) (*Participant, error) {
	rows, err := s.tx.Query(ctx, `
		select id, match_id, puuid, team_id, team_position, champion_id, champion_level, summoner_ids, rune_ids, item_ids, kills, deaths, assists, creep_score, damage_dealt, damage_taken, gold_earned, vision_score, pink_wards_bought
		from Participant
		where match_id = @match_id and puuid = @puuid;
	`, pgx.StrictNamedArgs{
		"match_id": matchID,
		"puuid":    puuid,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Participant])
}

func (s *ParticipantStore) CreateParticipant(ctx context.Context, req CreateParticipant) (*Participant, error) {
	rows, err := s.tx.Query(ctx, `
		insert into Participant (match_id, puuid, team_id, team_position, champion_id, champion_level, summoner_ids, rune_ids, item_ids, kills, deaths, assists, creep_score, damage_dealt, damage_taken, gold_earned, vision_score, pink_wards_bought)
		values (@match_id, @puuid, @team_id, @team_position, @champion_id, @champion_level, @summoner_ids, @rune_ids, @item_ids, @kills, @deaths, @assists, @creep_score, @damage_dealt, @damage_taken, @gold_earned, @vision_score, @pink_wards_bought)
		returning id, match_id, puuid, team_id, team_position, champion_id, champion_level, summoner_ids, rune_ids, item_ids, kills, deaths, assists, creep_score, damage_dealt, damage_taken, gold_earned, vision_score, pink_wards_bought;
	`, pgx.StrictNamedArgs{
		"match_id":          req.MatchID,
		"puuid":             req.PUUID,
		"team_id":           req.TeamID,
		"team_position":     req.TeamPosition,
		"champion_id":       req.ChampionID,
		"champion_level":    req.ChampionLevel,
		"summoner_ids":      req.SummonerIDs,
		"rune_ids":          req.RuneIDs,
		"item_ids":          req.ItemIDs,
		"kills":             req.Kills,
		"deaths":            req.Deaths,
		"assists":           req.Assists,
		"creep_score":       req.CreepScore,
		"damage_dealt":      req.DamageDealt,
		"damage_taken":      req.DamageTaken,
		"gold_earned":       req.GoldEarned,
		"vision_score":      req.VisionScore,
		"pink_wards_bought": req.PinkWardsBought,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Participant])
}

func (s *ParticipantStore) GetParticipantByMatchID(ctx context.Context, matchID string) ([]*Participant, error) {
	rows, err := s.tx.Query(ctx, `
		select id, match_id, puuid, team_id, team_position, champion_id, champion_level, summoner_ids, rune_ids, item_ids, kills, deaths, assists, creep_score, damage_dealt, damage_taken, gold_earned, vision_score, pink_wards_bought
		from Participant
		where match_id = @match_id;
	`, pgx.StrictNamedArgs{
		"match_id": matchID,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Participant])
}

type PageParam struct {
	Size   int
	LastID string
	Total  int
}

type OrderParam struct {
	Date bool
}

func (s *ParticipantStore) GetParticipantByPUUID(
	ctx context.Context,
	puuid string,
	page PageParam,
	order OrderParam,
) ([]*Participant, error) {
	rows, err := s.tx.Query(ctx, `
		select id, match_id, puuid, team_id, team_position, champion_id, champion_level, summoner_ids, rune_ids, item_ids, kills, deaths, assists, creep_score, damage_dealt, damage_taken, gold_earned, vision_score, pink_wards_bought
		from Participant
		order by date
		where puuid = @puuid and id > @last_id
		limit @size;
	`, pgx.StrictNamedArgs{
		"puuid": puuid,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Participant])
}
