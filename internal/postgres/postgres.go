package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Tx interface {
	Exec(ctx context.Context, sql string, args ...any) (commandTag pgconn.CommandTag, err error)

	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type BatchTx interface {
	Queue(query string, arguments ...any) *pgx.QueuedQuery
}

func errWrap(parent *error, format string, args ...any) {
	if *parent != nil {
		*parent = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), *parent)
	}
}
