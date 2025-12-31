package database

import (
	"context"
	"database/sql"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

// TxExec executes a DML query in a transaction.
// If no transaction is active, it falls back to Exec.
func (p *Postgres) TxExec(ctx context.Context, query string, args ...any) (int64, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if !ok {
		return p.Exec(ctx, query, args...)
	}

	var (
		tag pgconn.CommandTag
		err error
	)

	defer func() {
		if err != nil {
			tx.Rollback(ctx) //nolint:errcheck // nothing to do with it
		}
	}()

	tag, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// QueryRow executes a DQL query that must return at most one row.
func (p *Postgres) QueryRow(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Get(ctx, p.pool, dst, query, args...)
}

// TxQueryRow executes a DQL query that must return at most one row in a transaction.
// If no transaction is active, it falls back to QueryRow.
func (p *Postgres) TxQueryRow(ctx context.Context, dst any, query string, args ...any) error {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if !ok {
		return p.QueryRow(ctx, dst, query, args...)
	}

	var err error
	defer func() {
		if err != nil {
			tx.Rollback(ctx) //nolint:errcheck // nothing to do with it
		}
	}()
	err = pgxscan.Get(ctx, tx, dst, query, args...)
	return err
}

// QuerySlice executes a DQL query that returns multiple rows.
func (p *Postgres) QuerySlice(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Select(ctx, p.pool, dst, query, args...)
}

// TxQuerySlice executes a DQL query that returns multiple rows in a transaction.
// If no transaction is active, it falls back to QuerySlice.
func (p *Postgres) TxQuerySlice(ctx context.Context, dst any, query string, args ...any) error {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if !ok {
		return p.QuerySlice(ctx, dst, query, args...)
	}

	var err error
	defer func() {
		if err != nil {
			tx.Rollback(ctx) //nolint:errcheck // nothing to do with it
		}
	}()
	err = pgxscan.Select(ctx, tx, dst, query, args...)
	return err
}

// BeginTx starts a new transaction and returns a new context containing it.
func (p *Postgres) BeginTx(ctx context.Context, level TxLevel) (context.Context, error) {
	var opts pgx.TxOptions
	switch level {
	case TxLevelReadCommitted:
		opts.IsoLevel = pgx.ReadCommitted
	case TxLevelSerializable:
		opts.IsoLevel = pgx.Serializable
	}

	tx, err := p.pool.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, TxKey, tx)
	return ctx, nil
}

// CommitTx commits the transaction in the context.
func (p *Postgres) CommitTx(ctx context.Context) error {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if !ok {
		return ErrNoTx
	}
	return tx.Commit(ctx)
}

// RollbackTx rollbacks the transaction in the context.
func (p *Postgres) RollbackTx(ctx context.Context) error {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if !ok {
		return ErrNoTx
	}
	return tx.Rollback(ctx)
}
