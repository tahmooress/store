package psql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" //nolint: revive
	"github.com/tahmooress/store/internal/repository"
)

type DB struct {
	conn *sql.DB
}

type Config struct {
	DatabaseName        string
	DatabaseUser        string
	DatabasePassword    string
	DatabaseHost        string
	DatabasePort        string
	DatabaseMaxPageSize string
	DatabaseSSLMode     string
}

func New(ctx context.Context, config *Config) (*DB, error) {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DatabaseHost, config.DatabasePort, config.DatabaseUser,
		config.DatabasePassword, config.DatabaseName)

	conn, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, fmt.Errorf("psql opendb: %s", err)
	}

	if err = conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("psql ping: %s", err)
	}

	return &DB{
		conn: conn,
	}, nil
}

type transaction struct {
	tx *sql.Tx
}

var _ repository.Tx = &transaction{}

func (d *DB) beginTx(ctx context.Context) (*transaction, error) {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("beginTx: %s", err)
	}

	return &transaction{
		tx: tx,
	}, nil
}

func (d *DB) Exec(ctx context.Context, fn func(tx repository.Tx) error) error {
	tx, err := d.beginTx(ctx)
	if err != nil {
		return fmt.Errorf("psql Exec: %s", err)
	}

	err = fn(tx)
	if err != nil {
		_ = tx.tx.Rollback()

		return err
	}

	return tx.tx.Commit()
}

func (d *DB) Close() error {
	return d.conn.Close()
}
