// Package controller implements logic layer for balance handling.
package controller

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/balance/model"
)

type balanceRepo interface {
	InsertWithdrawal(ctx context.Context, userID int64, orderNumber string, amount float64) (int64, error)
	FindBalanceByUserID(ctx context.Context, userID int64) (*model.Balance, error)
	ListWithdrawals(ctx context.Context, userID int64) ([]*model.Withdrawal, error)
	UpdateBalanceWithdraw(ctx context.Context, userID int64, withdrawal float64) error
}

// Controller provides methods for balance handling logic.
type Controller struct {
	br balanceRepo
}

// New creates a new Controller instance.
func New(br balanceRepo) *Controller {
	return &Controller{br: br}
}
