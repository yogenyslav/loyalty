package controller

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/balance/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
	"github.com/yogenyslav/loyalty/pkg/luhn"
)

// Withdraw processes a withdrawal.
func (ctrl *Controller) Withdraw(ctx context.Context, userID int64, req *model.WithdrawReq) error {
	if !luhn.Validate(req.Order) {
		return errs.Wrap(luhn.ErrInvalidNumber, "invalid order number")
	}

	balance, err := ctrl.br.FindBalanceByUserID(ctx, userID)
	if err != nil {
		return errs.Wrap(err, "find balance by id")
	}
	if balance.Current < req.Sum {
		return errs.Wrap(errs.ErrNotEnoughBalance, "asked amount exceeds current balance")
	}

	_, err = ctrl.br.InsertWithdrawal(ctx, userID, req.Order, req.Sum)
	return errs.Wrap(err, "insert withdrawal")
}
