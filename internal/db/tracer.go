package db

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rank1zen/yujin/internal/logging"
	"go.uber.org/zap"
)

type tracer struct{}

func newTracer() pgx.QueryTracer {
	return &tracer{}
}

func (t *tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	logger := logging.FromContext(ctx).Sugar()

	q := strings.Join(strings.Fields(data.SQL), " ")
	logger.Debugf("query: %.20s", q)

	return ctx
}

func (t *tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	logger := logging.FromContext(ctx).Sugar()
	if data.Err != nil {
		logger.Debugf("%v: %v", data.CommandTag, data.Err)
	}
}

func (t *tracer) TraceBatchStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchStartData) context.Context {
	ctx = logging.WithFields(ctx, zap.String("mode", "batch"))
	logger := logging.FromContext(ctx).Sugar()

	logger.Debugf("%v queued", data.Batch.Len())

	return ctx
}

func (t *tracer) TraceBatchQuery(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchQueryData) {
	logger := logging.FromContext(ctx).Sugar()

	if data.Err != nil {
		sql := strings.Fields(data.SQL)
		logger.Debugf("%v: %v", sql[:4], data.Err)
	}
}

func (t *tracer) TraceBatchEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchEndData) {}
