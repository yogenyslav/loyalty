// Package database provides wrappers and utils for database interactions.
package database

import (
	"context"
	"database/sql"
)

// Config holds the database configuration settings.
type Config struct {
	DSN string `yaml:"dsn" env:"DB_DSN"`
}

// DB defines methods to operate with DB.
//
//go:generate mockgen -destination=../../tests/mocks/db.go -package=mocks . DB
type DB interface {
	Exec(ctx context.Context, query string, args ...any) (int64, error)
	QueryRow(ctx context.Context, dsy any, query string, args ...any) error
	QuerySlice(ctx context.Context, dst any, query string, args ...any) error
	Ping(ctx context.Context) error
	SQLDB() (*sql.DB, error)
	Close()
}
