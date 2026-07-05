package database

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pedro-h-castoldi/go-ddd-boilerplate/internal/config"
)

type IDatabase interface {
	Open(cfgs config.DatabaseConfig) error
	Close() error
	NewTransaction(ctx context.Context) (tx *sql.Tx, err error)
}

type IPool interface {
	OpenConnections(cfgs []config.DatabaseConfig) error
	NewTransaction(ctx context.Context, readonly bool, nickname string) (tx ITransactionManager, err error)
	Close() error
}

type ITransactionManager interface {
	Commit() error
	Rollback()
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Builder() sq.StatementBuilderType
}
