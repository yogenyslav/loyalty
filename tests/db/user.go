package db_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/loyalty/pkg/database"
)

const insertUser = `
	insert into loyalty.user (login, hashed_password)
	values ('test_user', 'hashed_password')
	returning id;
`

// InsertUser inserts a test user into the database and returns the user ID.
func InsertUser(t *testing.T, db database.DB) int64 {
	t.Helper()

	var userID int64
	err := db.QueryRow(t.Context(), &userID, insertUser)
	require.NoError(t, err, "insert user")

	return userID
}
