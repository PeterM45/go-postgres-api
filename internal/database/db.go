package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/PeterM45/go-postgres-api/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	Pool   *pgxpool.Pool
	Config *config.Config
}

func New(cfg *config.Config) (*DB, error) {
	pool, err := pgxpool.Connect(context.Background(), cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	// Create an instance with the pool
	db := &DB{Pool: pool, Config: cfg}

	// Initialize tables
	if err := db.createUserTable(); err != nil {
		pool.Close() // Close pool if table creation fails
		return nil, err
	}

	return db, nil
}

func (db *DB) createUserTable() error {
	fields := []string{
		fmt.Sprintf("id %s PRIMARY KEY", db.Config.User.IDField),
		"password_hash BYTEA NOT NULL",
		"created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP",
	}

	if db.Config.User.RequireUsername {
		fields = append(fields, "username VARCHAR(50) UNIQUE")
	}
	if db.Config.User.RequireEmail {
		fields = append(fields, "email VARCHAR(255) UNIQUE")
	}

	sql := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS users (
            %s
        )
    `, strings.Join(fields, ",\n            "))

	_, err := db.Pool.Exec(context.Background(), sql)
	return err
}

func (db *DB) Close() {
	db.Pool.Close()
}
