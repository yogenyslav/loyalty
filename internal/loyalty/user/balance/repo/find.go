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
