// Package repo provides methods to interact with orders in the database.
package repo

import (
	"context"

	"github.com/yogenyslav/loyalty/pkg/database"
)

// Repo is a repository for orders.
type Repo struct {
	db database.DB
}

// New creates a new Repo instance.
func New(db database.DB) *Repo {
	return &Repo{db: db}
}

// BeginTx starts a new database transaction.
func (r *Repo) BeginTx(ctx context.Context) (context.Context, error) {
	return r.db.BeginTx(ctx)
}

// CommitTx commits the current database transaction.
func (r *Repo) CommitTx(ctx context.Context) error {
	return r.db.CommitTx(ctx)
}

// RollbackTx rolls back the current database transaction.
func (r *Repo) RollbackTx(ctx context.Context) error {
	return r.db.RollbackTx(ctx)
}
