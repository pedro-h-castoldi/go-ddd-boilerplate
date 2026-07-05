package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pedro-h-castoldi/go-ddd-boilerplate/internal/config"
)

type pool struct {
	databases map[string]IDatabase
}

func NewPool(cfgs []config.DatabaseConfig) IPool {
	var pollMap = make(map[string]IDatabase)
	for _, cfg := range cfgs {
		pollMap[cfg.Nickname] = &postgres{}
	}
	return &pool{databases: pollMap}
}

func (p *pool) OpenConnections(cfgs []config.DatabaseConfig) error {
	for _, cfg := range cfgs {
		if _, ok := p.databases[cfg.Nickname]; ok {
			if err := p.databases[cfg.Nickname].Open(cfg); err != nil {
				return fmt.Errorf("open database %s: %w", cfg.Nickname, err)
			}
		} else {
			return fmt.Errorf("database %s not found", cfg.Nickname)
		}
	}
	return nil
}

func (p *pool) NewTransaction(ctx context.Context, readonly bool, nickname string) (ITransactionManager, error) {
	var (
		db  string = nickname
		tx  *sql.Tx
		err error
	)

	if db == "" {
		db = config.Get().GetMainDatabaseNickname()
	}

	if _, ok := p.databases[db]; ok {
		tx, err = p.databases[db].NewTransaction(ctx)
		if err != nil {
			return nil, fmt.Errorf("new transaction: %w", err)
		}
	} else {
		return nil, fmt.Errorf("database %s not found", db)
	}

	return NewTransactionManager(tx), nil
}

func (p *pool) Close() error {
	for _, db := range p.databases {
		if err := db.Close(); err != nil {
			return fmt.Errorf("close database: %w", err)
		}
	}
	return nil
}
