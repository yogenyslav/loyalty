package database

import (
	"context"

	"github.com/yogenyslav/loyalty/pkg/errs"
)

// UnitOfWork provides methods to execute operations within a transaction.
//
//go:generate mockgen -destination=../../tests/mocks/uow.go -package=mocks . UnitOfWork
type UnitOfWork interface {
	// WithTx executes the given function within a database transaction.
	WithTx(ctx context.Context, level TxLevel, fn func(ctx context.Context) error) error
}

type unitOfWork struct {
	db DB
}

// NewUnitOfWork creates a new UnitOfWork instance.
func NewUnitOfWork(db DB) *unitOfWork {
	return &unitOfWork{db: db}
}

// WithTx implements UnitOfWork method.
func (uow *unitOfWork) WithTx(ctx context.Context, level TxLevel, fn func(ctx context.Context) error) error {
	tx, err := uow.db.beginTx(ctx, level)
	if err != nil {
		return errs.Wrap(err, "begin transaction")
	}

	defer func() {
		if e := recover(); e != nil {
			uow.db.rollbackTx(tx) //nolint:errcheck // nothing we can do
			panic(e)
		}

		if err != nil {
			uow.db.rollbackTx(tx) //nolint:errcheck // nothing we can do
		} else {
			err = uow.db.commitTx(tx)
		}
	}()

	err = fn(tx)
	return err
}
