package repo

import (
	"context"

	"github.com/yogenyslav/loyalty/pkg/errs"
)

const insertBalance = `
	insert into loyalty.balance (fk_user_id)
	values ($1);
`

// InsertBalance creates a new balance record for the specified user.
func (r *Repo) InsertBalance(ctx context.Context, userID int64) error {
	_, err := r.db.TxExec(ctx, insertBalance, userID)
	return errs.Wrap(err, "exec query")
}

const insertWithdrawal = `
	insert into loyalty.withdrawal (fk_user_id, order_number, amount)
	values ($1, $2, $3)
	returning id;
`

// InsertWithdrawal inserts a withdrawal record for specified user and order.
func (r *Repo) InsertWithdrawal(ctx context.Context, userID int64, orderNumber string, amount float64) (int64, error) {
	var withdrawalID int64
	err := r.db.TxQueryRow(ctx, &withdrawalID, insertWithdrawal, userID, orderNumber, amount)
	if err != nil {
		return 0, errs.Wrap(err, "tx query row")
	}
	return withdrawalID, nil
}
