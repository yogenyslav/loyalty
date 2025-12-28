// Package repo provides methods to interact with orders in the database.
package repo

import "github.com/yogenyslav/loyalty/pkg/database"

// Repo is a repository for orders.
type Repo struct {
	db database.DB
}

// New creates a new Repo instance.
func New(db database.DB) *Repo {
	return &Repo{db: db}
}
