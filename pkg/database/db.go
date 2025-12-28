// Package database provides wrappers and utils for database interactions.
package database

import (
	"context"
	"database/sql"
	"errors"
)

// CtxKey is a wrapper type for context keys.
type CtxKey string

// TxKey is the context key for database transactions.
const TxKey CtxKey = "tx"

var ErrNoTx = errors.New("no transaction in context")

// Config holds the database configuration settings.
type Config struct {
	URI      string `yaml:"uri"      env:"DATABASE_URI"`
	User     string `yaml:"user"     env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Name     string `yaml:"name"     env:"DB_NAME"`
	Host     string `yaml:"host"     env:"DB_HOST"`
	Port     string `yaml:"port"     env:"DB_PORT"`
	Driver   string `yaml:"driver"   env:"DB_DRIVER"`
	SSLMode  string `yaml:"sslmode"  env:"DB_SSLMODE"`
}

// DSN returns the Data Source Name for connecting to the database.
func (c *Config) DSN() string {
	if c.URI != "" {
		return c.URI
	}
	return c.Driver + "://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + c.Name + "?sslmode=" + c.SSLMode
}

// DB defines methods to operate with DB.
//
//go:generate mockgen -destination=../../tests/mocks/db.go -package=mocks . DB
type DB interface {
	Exec(ctx context.Context, query string, args ...any) (int64, error)
	TxExec(ctx context.Context, query string, args ...any) (int64, error)
	QueryRow(ctx context.Context, dst any, query string, args ...any) error
	TxQueryRow(ctx context.Context, dst any, query string, args ...any) error
	QuerySlice(ctx context.Context, dst any, query string, args ...any) error
	TxQuerySlice(ctx context.Context, dst any, query string, args ...any) error
	Ping(ctx context.Context) error
	SQLDB() (*sql.DB, error)
	Close()
	BeginTx(ctx context.Context) (context.Context, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error
}
