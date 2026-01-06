package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	t.Parallel()

	t.Run("Wrap error with context", func(t *testing.T) {
		rawErr := errors.New("original error")
		wrappedErr := Wrap(rawErr, "additional context")

		assert.ErrorIs(t, wrappedErr, rawErr)
		assert.Contains(t, wrappedErr.Error(), "additional context")
		assert.Contains(t, wrappedErr.Error(), "original error")
	})

	t.Run("Wrap nil error returns nil", func(t *testing.T) {
		var rawErr error
		wrappedErr := Wrap(rawErr, "additional context")

		assert.Nil(t, wrappedErr)
	})

	t.Run("Multiple wraps preserve original error", func(t *testing.T) {
		rawErr := errors.New("original error")
		firstWrap := Wrap(rawErr, "first context")
		secondWrap := Wrap(firstWrap, "second context")

		assert.ErrorIs(t, secondWrap, rawErr)
		assert.Contains(t, secondWrap.Error(), "first context")
		assert.Contains(t, secondWrap.Error(), "second context")
		assert.Contains(t, secondWrap.Error(), "original error")
	})

	t.Run("Wrap with empty context", func(t *testing.T) {
		rawErr := errors.New("original error")
		wrappedErr := Wrap(rawErr, "")

		assert.ErrorIs(t, wrappedErr, rawErr)
	})
}
