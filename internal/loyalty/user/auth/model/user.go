// Package model contains data models related to authentication.
package model

import "time"

// User is a data model for user in database.
type User struct {
	ID             int64     `db:"id"`
	Login          string    `db:"login"`
	HashedPassword string    `db:"hashed_password"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// UserDto is a data transfer object for user.
type UserDto struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}
