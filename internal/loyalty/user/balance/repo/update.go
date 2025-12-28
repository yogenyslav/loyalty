package repo

import (
	"context"

	"github.com/yogenyslav/loyalty/pkg/errs"
)

const updateBalance = `
	update loyalty.balance
	set current = current - $2,
	    withdrawn = withdrawn + $2,
	    updated_at = current_timestamp
	where fk_user_id = $1;
`

// UpdateBalanceWithdraw updates the balance for the specified user by subtracting withdrawal.
func (r *Repo) UpdateBalanceWithdraw(ctx context.Context, userID int64, withdrawal float64) error {
	_, err := r.db.TxExec(ctx, updateBalance, userID, withdrawal)
	return errs.Wrap(err, "exec query")
}

const updateBalanceAccrual = `
	update loyalty.balance
	set current = current + $2,
	    updated_at = current_timestamp
	where fk_user_id = $1;
`

// UpdateBalanceAccrual updates the balance for the specified user by adding accrual.
func (r *Repo) UpdateBalanceAccrual(ctx context.Context, userID int64, accrual float64) error {
	_, err := r.db.TxExec(ctx, updateBalanceAccrual, userID, accrual)
	return errs.Wrap(err, "exec query")
}
