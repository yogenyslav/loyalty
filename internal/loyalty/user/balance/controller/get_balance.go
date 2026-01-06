package controller

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/balance/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

// GetBalance returns user's balance.
func (ctrl *Controller) GetBalance(ctx context.Context, userID int64) (*model.BalanceDto, error) {
	balance, err := ctrl.br.FindBalanceByUserID(ctx, userID)
	if err != nil {
		return nil, errs.Wrap(err, "find balance by id")
	}
	return balance.ToDto(), nil
}
