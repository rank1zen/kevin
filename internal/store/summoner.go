package store

import (
	"context"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
)

type Summoner struct {
	ID            uuid.UUID `db:"id"`
	PUUID         string    `db:"puuid"`
	Region        string    `db:"region"`
	Name          string    `db:"name"`
	Tagline       string    `db:"tagline"`
	SummonerLevel int       `db:"summoner_level"`
	ProfileIconID string    `db:"profile_icon_id"`
	LastUpdated   time.Time `db:"last_updated"`
}

type CreateSummoner struct {
	PUUID         string
	Region        string
	Name          string
	Tagline       string
	SummonerLevel int
	ProfileIconID string
	LastUpdated   time.Time
}

type UpdateSummoner struct {
	Region        string
	Name          string
	Tagline       string
	SummonerLevel int
	ProfileIconID string
	LastUpdated   time.Time
}

type SummonerStore struct {
	tx DBTX
}

func NewSummonerStore(tx DBTX) *SummonerStore {
	return &SummonerStore{
		tx: tx,
	}
}

func (s *SummonerStore) GetSummonerByPUUID(ctx context.Context, puuid string) (*Summoner, error) {
	rows, err := s.tx.Query(ctx, `
		select id, puuid, name, tagline, region, summoner_level, profile_icon_id, last_updated
		from Summoner
		where puuid = @puuid;
	`, pgx.StrictNamedArgs{
		"puuid": puuid,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Summoner])
}

func (s *SummonerStore) CreateSummoner(ctx context.Context, req CreateSummoner) (*Summoner, error) {
	rows, err := s.tx.Query(ctx, `
		insert into Summoner (puuid, name, tagline, region, summoner_level, profile_icon_id, last_updated)
		values (@puuid, @name, @tagline, @region, @summoner_level, @profile_icon_id, @last_updated)
		returning id, puuid, name, tagline, region, summoner_level, profile_icon_id, last_updated;
	`, pgx.StrictNamedArgs{
		"puuid":           req.PUUID,
		"region":          req.Region,
		"name":            req.Name,
		"tagline":         req.Tagline,
		"summoner_level":  req.SummonerLevel,
		"profile_icon_id": req.ProfileIconID,
		"last_updated":    req.LastUpdated,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Summoner])
}

func (s *SummonerStore) UpdateSummonerByPUUID(ctx context.Context, puuid string, req UpdateSummoner) (*Summoner, error) {
	rows, err := s.tx.Query(ctx, `
		update Summoner
		set region = @region, name = @name, tagline = @tagline, summoner_level = @summoner_level, profile_icon_id = @profile_icon_id, last_updated = @last_updated
		where puuid = @puuid
		returning id, region, puuid, name, tagline, summoner_level, profile_icon_id, last_updated;
	`, pgx.StrictNamedArgs{
		"puuid":           puuid,
		"region":          req.Region,
		"name":            req.Name,
		"tagline":         req.Tagline,
		"summoner_level":  req.SummonerLevel,
		"profile_icon_id": req.ProfileIconID,
		"last_updated":    req.LastUpdated,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Summoner])
}

func (s *SummonerStore) GetSummonerByRegionNameTag(ctx context.Context, region, name, tag string) (*Summoner, error) {
	rows, err := s.tx.Query(ctx, `
		select id, puuid, region, name, tagline, summoner_level, profile_icon_id, last_updated
		from Summoner
		where region = @region and name = @name and tagline = @tag
	`, pgx.StrictNamedArgs{
		"region": region,
		"name":   name,
		"tag":    tag,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Summoner])
}

func (s *SummonerStore) SearchSummoner(ctx context.Context, query string) ([]*Summoner, error) {
	rows, err := s.tx.Query(ctx, `
		select id, puuid, region, name, tagline, summoner_level, profile_icon_id, last_updated
		from Summoner
		where name % @query or tagline % @query
		order by greatest(similarity(name, @query), similarity(tagline, @query) * 0.5) desc
		limit 20
	`, pgx.StrictNamedArgs{
		"query": query,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Summoner])
}
