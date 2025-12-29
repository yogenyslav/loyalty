package controller

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/balance/model"
)

// ListWithdrawals returns a list of user's withdrawals.
func (ctrl *Controller) ListWithdrawals(ctx context.Context, userID int64) ([]*model.WithdrawalDto, error) {
	withdrawals, err := ctrl.br.ListWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}

	var withdrawalsDto []*model.WithdrawalDto
	for _, w := range withdrawals {
		withdrawalsDto = append(withdrawalsDto, w.ToDto())
	}
	return withdrawalsDto, nil
}
