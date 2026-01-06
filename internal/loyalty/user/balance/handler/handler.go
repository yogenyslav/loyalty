// Package handler routes balance-related HTTP requests.
package handler

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/balance/model"
)

type balanceController interface {
	GetBalance(ctx context.Context, userID int64) (*model.BalanceDto, error)
	Withdraw(ctx context.Context, userID int64, req *model.WithdrawReq) error
	ListWithdrawals(ctx context.Context, userID int64) ([]*model.WithdrawalDto, error)
}

// Handler provides methods for handling balance-related HTTP requests.
type Handler struct {
	ctrl balanceController
}

// New creates a new Handler instance.
func New(ctrl balanceController) *Handler {
	return &Handler{ctrl: ctrl}
}
