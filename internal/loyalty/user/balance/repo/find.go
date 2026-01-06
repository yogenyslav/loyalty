package repo

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/balance/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

const findBalanceByUserID = `
	select fk_user_id, current, withdrawn, updated_at
	from loyalty.balance
	where fk_user_id = $1
`

// FindBalanceByUserID returns the balance for the given user.
func (r *Repo) FindBalanceByUserID(ctx context.Context, userID int64) (*model.Balance, error) {
	var balance model.Balance
	err := r.db.QueryRow(ctx, &balance, findBalanceByUserID, userID)
	if err != nil {
		return nil, errs.Wrap(err, "query row")
	}
	return &balance, nil
}

const listWithdrawals = `
	select id, amount, order_number, fk_user_id, created_at
	from loyalty.withdrawal
	where fk_user_id = $1
	order by created_at;
`

// ListWithdrawals returns the list of withdrawals ordered by creation date.
func (r *Repo) ListWithdrawals(ctx context.Context, userID int64) ([]*model.Withdrawal, error) {
	var withdrawals []*model.Withdrawal
	err := r.db.QuerySlice(ctx, &withdrawals, listWithdrawals, userID)
	return withdrawals, errs.Wrap(err, "query slice")
}
