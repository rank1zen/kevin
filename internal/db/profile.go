package db

import (
	"context"
	"fmt"

	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/pgxutil"
)

func (db *DB) CheckProfileExists(ctx context.Context, puuid internal.PUUID) (bool, error) {
	var exists bool
	err := db.pool.QueryRow(ctx, `
	SELECT EXISTS (
		SELECT 1
		FROM
			profiles
		WHERE
			puuid = $1
	);`, puuid).Scan(&exists)

	return exists, err
}

func upsertProfile(ctx context.Context, conn pgxutil.Exec, profile internal.Profile) error {
	_, err := conn.Exec(ctx, `
	INSERT INTO profiles
		(puuid, name, tagline, last_updated)
	VALUES
		($1, $2, $3, $4)
	ON CONFLICT
		(puuid)
	DO UPDATE SET
		name = $2,
		tagline = $3,
		last_updated = $4;
	`, profile.Puuid, profile.Name, profile.Tagline, profile.RecordDate)

	return err
}

func (db *DB) UpdateProfile(ctx context.Context, profile internal.Profile) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = createSummonerRecord(ctx, tx, profile)
	if err != nil {
		return fmt.Errorf("inserting summoner: %w", err)
	}

	err = createLeagueRecord(ctx, tx, profile)
	if err != nil {
		return fmt.Errorf("inserting league: %w", err)
	}

	err = upsertProfile(ctx, tx, profile)
	if err != nil {
		return fmt.Errorf("upserting profile: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetProfile(ctx context.Context, puuid internal.PUUID) (internal.Profile, error) {
	var p internal.Profile
	err := db.pool.QueryRow(ctx, `
	SELECT
		puuid,
		last_updated,
		name,
		tagline
	FROM
		profiles
	WHERE
		puuid = $1;
	`, puuid).Scan(
		&p.Puuid,
		&p.AccountID,
		&p.Level,
		&p.Name,
		&p.ProfileIconID,
		&p.RecordDate,
		&p.RevisionDate,
		&p.SummonerID,
		&p.Tagline,
		&p.ValidFrom,
		&p.ValidTo,
	)

	return p, err
}
