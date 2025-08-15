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
			"puuid": summoner.PUUID,
			"name": summoner.Name,
			"tagline": summoner.Tagline,
		},
	)

	return err
}

func (db *SummonerStore) GetSummoner(ctx context.Context, puuid riot.PUUID) (Summoner, error) {
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
