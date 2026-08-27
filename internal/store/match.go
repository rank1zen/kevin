package store

import (
	"context"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
)

type Match struct {
	ID       uuid.UUID     `db:"id"`
	Region   string        `db:"region"`
	MatchID  string        `db:"match_id"`
	Version  string        `db:"version"`
	Date     time.Time     `db:"date"`
	Duration time.Duration `db:"duration"`
	WinnerID string        `db:"winner_id"`
}

type CreateMatch struct {
	Region   string
	MatchID  string
	Version  string
	Date     time.Time
	Duration time.Duration
	WinnerID string
}

type GetMatchByPUUIDPageParams struct {
	Ascending bool

	Size     int
	LastID   uuid.UUID
	LastDate time.Time
}

type MatchStore struct {
	tx DBTX
}

func NewMatchStore(tx DBTX) *MatchStore {
	return &MatchStore{
		tx: tx,
	}
}

func (s *MatchStore) GetMatchByMatchID(ctx context.Context, matchID string) (*Match, error) {
	rows, err := s.tx.Query(ctx, `
		select id, region, match_id, version, date, duration, winner_id
		from Match
		where match_id = @match_id;
	`, pgx.StrictNamedArgs{
		"match_id": matchID,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Match])
}

func (s *MatchStore) CreateMatch(ctx context.Context, req CreateMatch) (*Match, error) {
	rows, err := s.tx.Query(ctx, `
		insert into Match (region, match_id, version, date, duration, winner_id)
		values (@region, @match_id, @version, @date, @duration, @winner_id)
		returning id, region, match_id, version, date, duration, winner_id;
	`, pgx.StrictNamedArgs{
		"region":    req.Region,
		"match_id":  req.MatchID,
		"version":   req.Version,
		"date":      req.Date,
		"duration":  req.Duration,
		"winner_id": req.WinnerID,
	})
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Match])
}

// GetMatchByPUUID finds matches where the there is a participant with the given puuid.
func (s *MatchStore) GetMatchByPUUID(
	ctx context.Context,
	puuid string,
	req GetMatchByPUUIDPageParams,
) (
	[]*Match,
	GetMatchByPUUIDPageParams,
	error,
) {
	var nextPage GetMatchByPUUIDPageParams

	var pageSize = 10
	if req.Size > 0 {
		pageSize = req.Size
	}

	var sql = `
		select m.id, m.region, m.match_id, m.version, m.date, m.duration, m.winner_id
		from Match m
		join Participant p on m.match_id = p.match_id
	`

	args := pgx.StrictNamedArgs{}

	// SQL string building

	sql += ` where p.puuid = @puuid`
	args["puuid"] = puuid

	if req.LastID != uuid.Nil() && !req.LastDate.IsZero() {
		if req.Ascending {
			sql += ` and (@last_date, @last_id) < (m.date, m.id)`
		} else {
			sql += ` and (@last_date, @last_id) > (m.date, m.id)`
		}
		args["last_id"] = req.LastID
		args["last_date"] = req.LastDate
	}

	if req.Ascending {
		sql += ` order by m.date asc`
	} else {
		sql += ` order by m.date desc`
	}
	nextPage.Ascending = req.Ascending

	sql += ` limit @size;`
	args["size"] = pageSize
	nextPage.Size = pageSize

	rows, err := s.tx.Query(ctx, sql, args)
	if err != nil {
		return nil, GetMatchByPUUIDPageParams{}, err
	}

	results, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Match])
	if err != nil {
		return nil, GetMatchByPUUIDPageParams{}, err
	}

	if len(results) >= pageSize {
		nextPage.LastID = results[len(results)-1].ID
		nextPage.LastDate = results[len(results)-1].Date
	} else {
		nextPage.LastID = uuid.Nil()
		nextPage.LastDate = time.Time{}
	}

	return results, nextPage, nil
}
