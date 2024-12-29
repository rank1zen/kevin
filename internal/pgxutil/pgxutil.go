package pgxutil

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SqlValue string

type Conn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type Query interface {
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row
}

type Exec interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func sanitizeIdentifier(s string) string {
	return pgx.Identifier{s}.Sanitize()
}

func makePgxIdentifier(v any) (pgx.Identifier, error) {
	switch v := v.(type) {
	case string:
		return pgx.Identifier{v}, nil
	case pgx.Identifier:
		return v, nil
	default:
		return nil, fmt.Errorf("expected string or pgx.Identifier, got %T", v)
	}
}

func QueryInsertRow(ctx context.Context, db Query, tableName any, values map[string]any) error {
	tableIdent, err := makePgxIdentifier(tableName)
	if err != nil {
		return fmt.Errorf("InsertRow invalid tableName: %w", err)
	}

	sql, args := buildInsertRowSql(tableIdent, values, "")
	_, err = wrapExec(ctx, db, sql, args)
	return err
}

func BatchInsertRow(batch *pgx.Batch, tableName any, values map[string]any) {
	tableIdent, err := makePgxIdentifier(tableName)
	if err != nil {
		// Panicking is undesirable, but we don't want to have this function return an error or silently ignore the error.
		// Possibly pgx.Batch should be modified to allow queueing an error.
		panic(fmt.Sprintf("batchInsertRow invalid tableName: %v", err))
	}

	sql, args := buildInsertRowSql(tableIdent, values, "")
	batch.Queue(sql, args...)
}

func buildInsertRowSql(tableName pgx.Identifier, values map[string]any, returningClause string) (sql string, args []any) {
	b := &strings.Builder{}
	b.WriteString("insert into ")
	if len(tableName) == 1 {
		b.WriteString(sanitizeIdentifier(tableName[0]))
	} else {
		b.WriteString(tableName.Sanitize())
	}
	b.WriteString(" (")

	// Go maps are iterated in random order. The generated SQL should be stable so sort the keys.
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		if i > 0 {
			b.WriteString(", ")
		}
		sanitizedKey := sanitizeIdentifier(k)
		b.WriteString(sanitizedKey)
	}

	b.WriteString(") values (")
	args = make([]any, 0, len(keys))
	for _, k := range keys {
		if len(args) > 0 {
			b.WriteString(", ")
		}
		if SQLValue, ok := values[k].(SqlValue); ok {
			b.WriteString(string(SQLValue))
		} else {
			args = append(args, values[k])
			b.WriteByte('$')
			b.WriteString(strconv.FormatInt(int64(len(args)), 10))
		}
	}

	b.WriteString(")")

	if returningClause != "" {
		b.WriteString(" returning ")
		b.WriteString(returningClause)
	}

	return b.String(), args
}

func wrapExec(ctx context.Context, db Query, sql string, args []any) (pgconn.CommandTag, error) {
	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	rows.Close()

	return rows.CommandTag(), rows.Err()
}
