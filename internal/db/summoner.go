package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/pgxutil"
)

// TODO: fix these types
type summonerRecordRow struct {
	ValidFrom     bool
	ValidTo       bool
	EnteredAt     bool
	AccountId     bool
	SummonerId    bool
	Puuid         bool
	RevisionDate  bool
	SummonerLevel bool
	ProfileIconID bool
}

func createSummonerRecord(ctx context.Context, conn pgxutil.Query, m internal.Profile) (summonerRecordRow, error) {
	row := pgx.NamedArgs{
		"valid_from":      1, // TODO
		"valid_to":        1,
		"entered_at":      m.RecordDate,
		"account_id":      m.AccountID,
		"summoner_id":     m.SummonerID,
		"puuid":           m.Puuid,
		"revision_date":   m.RevisionDate,
		"summoner_level":  m.Level,
		"profile_icon_id": m.ProfileIconID,
	}

	var inserted summonerRecordRow
	err := conn.QueryRow(ctx, `
	INSERT INTO summoner_records (
		valid_from,
		valid_to,
		entered_at,
		account_id,
		summoner_id,
		puuid,
		revision_date,
		summoner_level,
		profile_icon_id
	)
	VALUES (
		@valid_from,
		@valid_to,
		@entered_at,
		@account_id,
		@summoner_id,
		@puuid,
		@revision_date,
		@summoner_level,
		@profile_icon_id
	)
	RETURNING
		valid_from,
		valid_to,
		entered_at,
		account_id,
		summoner_id,
		puuid,
		revision_date,
		summoner_level,
		profile_icon_id;
	`, row).Scan(
		&inserted.ValidFrom,
		&inserted.ValidTo,
		&inserted.EnteredAt,
		&inserted.AccountId,
		&inserted.SummonerId,
		&inserted.Puuid,
		&inserted.RevisionDate,
		&inserted.SummonerLevel,
		&inserted.ProfileIconID,
	)
	if err != nil {
		return summonerRecordRow{}, err
	}

	return inserted, nil
}

func getSummonerRecord(ctx context.Context, conn pgxutil.Query, id string) error {
	conn.Query(ctx, `
	SELECT
		valid_from,
		valid_to,
		entered_at,
		account_id,
		summoner_id,
		puuid,
		revision_date,
		summoner_level,
		profile_icon_id
	FROM
	`)

	return nil
}
