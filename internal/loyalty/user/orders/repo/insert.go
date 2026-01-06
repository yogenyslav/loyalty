package repo

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

const insertOrder = `
	insert into loyalty.order(number, accrual, fk_user_id)
	values ($1, $2, $3);
`

// InsertOrder inserts a new order into the database.
func (r *Repo) InsertOrder(ctx context.Context, order *model.Order) error {
	_, err := r.db.Exec(ctx, insertOrder, order.Number, order.Accrual, order.UserID)
	return errs.Wrap(err, "exec query")
}

const shedulePolling = `
	insert into loyalty.poll_accrual(order_number)
	values ($1);
`

// SchedulePolling adds an order to the polling schedule.
func (r *Repo) SchedulePolling(ctx context.Context, orderNumber string) error {
	_, err := r.db.Exec(ctx, shedulePolling, orderNumber)
	return errs.Wrap(err, "exec query")
}
