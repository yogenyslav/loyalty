// Package repo provides methods to interact with user authentication data in the database.
package repo

import (
	"context"

	"github.com/yogenyslav/loyalty/pkg/database"
)

type balanceRepo interface {
	InsertBalance(ctx context.Context, userID int64) error
}

// Repo is a repository for user authentication data.
type Repo struct {
	db database.DB
	br balanceRepo
}

// New creates a new Repo instance.
func New(db database.DB, br balanceRepo) *Repo {
	return &Repo{
		db: db,
		br: br,
	}
}
