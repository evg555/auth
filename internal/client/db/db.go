package db

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Client interface {
	DB() DB
	Close() error
}

type Query struct {
	Name     string
	QueryRaw string
}

type SQLExecutor interface {
	NamedExecutor
	QueryExecutor
}

type NamedExecutor interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

type QueryExecutor interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRawContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type DB interface {
	SQLExecutor
	Pinger
	Close()
}
