package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/kevin/internal/riot"
)

type Summoner struct {
	PUUID riot.PUUID

	Name, Tagline string
}

// SummonerStore manages summoner objects.
type SummonerStore struct{ Tx Tx }

func (db *SummonerStore) CreateSummoner(ctx context.Context, summoner Summoner) error {
	_, err := db.Tx.Exec(ctx, `
		INSERT INTO Summoner (
			puuid,
			name,
			tagline
		)
		VALUES (
			@puuid,
			@name,
			@tagline
		)
		ON CONFLICT (puuid)
		DO UPDATE SET
			name    = @name,
			tagline = @tagline;
	`,
		pgx.NamedArgs{
			"puuid":   summoner.PUUID,
			"name":    summoner.Name,
			"tagline": summoner.Tagline,
		},
	)

	return err
}

func (db *SummonerStore) GetSummoner(ctx context.Context, puuid riot.PUUID) (_ Summoner, err error) {
	defer errWrap(&err, "SummonerStore.GetSummoner(ctx, %v)", puuid)

	rows, err := db.Tx.Query(ctx, `
		SELECT
			puuid,
			name,
			tagline
		FROM
			Summoner
		WHERE
			puuid = $1;
	`, puuid)

	if err != nil {
		return Summoner{}, err
	}

	summoner, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Summoner])
	if err != nil {
		return Summoner{}, err
	}

	return summoner, nil
}

// TODO: switch to fuzzy match algorithm.
func (db *SummonerStore) SearchSummoner(ctx context.Context, q string) ([]Summoner, error) {
	rows, err := db.Tx.Query(ctx, `
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

	if err != nil {
		return nil, err
	}

	summonerResults, err := pgx.CollectRows(rows, pgx.RowToStructByName[Summoner])
	if err != nil {
		return nil, err
	}

	return summonerResults, nil
}

// SearchByNameTag searches for summoners by name and tag. Names will be matched
// by prefix. If tag is empty, search will not match by tag. If tag is given,
// search will match by tag exactly. Search is case-insensitive and returns 10
// results.
//
// NOTE: The search algorithm is up for change.
func (db *SummonerStore) SearchByNameTag(ctx context.Context, name, tag string) ([]Summoner, error) {
	q := `
		SELECT
			puuid,
			name,
			tagline
		FROM
			Summoner
		WHERE
			name LIKE @name || '%'
	`

	args := pgx.NamedArgs{
		"name": name,
	}

	if tag != "" {
		q += `
			AND tagline = @tag
		`
		args["tag"] = tag
	}

	q += `
		LIMIT 10;
	`

	rows, err := db.Tx.Query(ctx, q, args)
	if err != nil {
		return nil, err
	}

	summonerResults, err := pgx.CollectRows(rows, pgx.RowToStructByName[Summoner])
	if err != nil {
		return nil, err
	}

	return summonerResults, nil
}
