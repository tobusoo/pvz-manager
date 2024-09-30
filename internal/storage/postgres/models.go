package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type QueryEngine interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TransactionManager interface {
	GetQueryEngine(ctx context.Context) QueryEngine
	RunReadOnlyCommitted(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunReadUncommitted(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunReadCommitted(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error
}
