package errs

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// 400 Bad Request
var (
	// ErrInvalidRequest is an error when invalid request is provided.
	ErrInvalidRequest = errors.New("invalid request")
)

// 401 Unauthorized
var (
	// ErrInvalidCredentials is an error when provided credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// 409 Conflict
var (
	// ErrUserAlreadyExists is an error when trying to create a user with login that already exists.
	ErrUserAlreadyExists = errors.New("user with such login already exists")
)

// CheckUniqueViloation checks if the error is a postgres unique violation error.
func CheckUniqueViloation(err error) bool {
	var pgError *pgconn.PgError
	return errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation
}
