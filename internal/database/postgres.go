package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pedro-h-castoldi/go-ddd-boilerplate/internal/config"
)

type postgres struct {
	db                 *sql.DB
	transactionTimeout time.Duration
}

func (p *postgres) Open(cfg config.DatabaseConfig) error {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)

	if cfg.SSLMode == "require" {
		dsn += fmt.Sprintf(" sslcert=%s sslkey=%s sslrootcert=%s", cfg.SSLConfig.ClientCert, cfg.SSLConfig.ClientKey, cfg.SSLConfig.RootCA)
	}

	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("parse postgres config: %w", err)
	}

	db := stdlib.OpenDB(*pgxConfig.ConnConfig)
	if err != nil {
		return fmt.Errorf("open postgres connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(1 * time.Minute)

	p.db = db
	p.transactionTimeout = time.Duration(cfg.TransactionTimeout) * time.Second

	return nil
}

func (p *postgres) Close() error {
	return p.db.Close()
}

func (p *postgres) NewTransaction(ctx context.Context) (tx *sql.Tx, err error) {

	ctx, cancel := context.WithTimeout(ctx, p.transactionTimeout)
	defer cancel()

	tx, err = p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	return tx, nil
}
