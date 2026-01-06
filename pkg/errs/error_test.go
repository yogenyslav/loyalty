package errs

import (
	"errors"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestCheckUniqueViloation(t *testing.T) {
	t.Parallel()

	t.Run("unique violation error", func(t *testing.T) {
		t.Parallel()

		err := &pgconn.PgError{Code: pgerrcode.UniqueViolation}
		isUniqueViolation := CheckUniqueViloation(err)
		assert.True(t, isUniqueViolation)
	})

	t.Run("other pg error", func(t *testing.T) {
		t.Parallel()

		err := &pgconn.PgError{Code: pgerrcode.ForeignKeyViolation}
		isUniqueViolation := CheckUniqueViloation(err)
		assert.False(t, isUniqueViolation)
	})

	t.Run("non-pg error", func(t *testing.T) {
		t.Parallel()

		err := errors.New("error")
		isUniqueViolation := CheckUniqueViloation(err)
		assert.False(t, isUniqueViolation)
	})
}
