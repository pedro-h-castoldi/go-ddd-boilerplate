package database

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type transactionManager struct {
	tx      *sql.Tx
	builder sq.StatementBuilderType
}

func NewTransactionManager(tx *sql.Tx) ITransactionManager {
	return &transactionManager{tx: tx, builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(tx)}
}

func (tm *transactionManager) Commit() error {
	return tm.tx.Commit()
}

func (tm *transactionManager) Rollback() {
	_ = tm.tx.Rollback()
}

func (tm *transactionManager) Query(query string, args ...any) (*sql.Rows, error) {
	return tm.tx.Query(query, args...)
}

func (tm *transactionManager) QueryRow(query string, args ...any) *sql.Row {
	return tm.tx.QueryRow(query, args...)
}

func (tm *transactionManager) Builder() sq.StatementBuilderType {
	return tm.builder
}
