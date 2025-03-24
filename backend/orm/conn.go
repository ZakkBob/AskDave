package orm

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pashagolub/pgxmock/v4"
)

type PgxPoolIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error)
	Close()
}

var dbpool PgxPoolIface

func MockConnect(mock pgxmock.PgxPoolIface) {
	dbpool = mock
}

func Connect(url string) error {
	if url == "" {
		url = "postgres://postgres:password@127.0.0.1:5432/postgres"
	}
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	dbpool = pool
	return nil
}

func Close() {
	dbpool.Close()
}

func Ping(c context.Context) error {
	return dbpool.Ping(c)
}

func DbPool() PgxPoolIface {
	return dbpool
}
