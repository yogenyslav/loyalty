package database

import (
	"context"
	"database/sql"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres is a wrapper around PostgreSQL using pgx connection pool.
type Postgres struct {
	pool *pgxpool.Pool
}

// NewPostgres creates a new pg instance.
func NewPostgres(ctx context.Context, dsn string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		pool: pool,
	}, nil
}

// SQLDB return a sql.DB format database conn.
func (p *Postgres) SQLDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", p.pool.Config().ConnString())
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Close the underlying pool.
func (p *Postgres) Close() {
	p.pool.Close()
}

// Ping the database.
func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

// Exec executes a DML query.
func (p *Postgres) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	tag, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// QueryRow executes a DQL query that must return at most one row.
func (p *Postgres) QueryRow(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Get(ctx, p.pool, dst, query, args...)
}

// QuerySlice executes a DQL query that returns multiple rows.
func (p *Postgres) QuerySlice(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Select(ctx, p.pool, dst, query, args...)
}
