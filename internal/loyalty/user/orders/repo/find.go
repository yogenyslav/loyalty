package repo

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

const listOrders = `
	select number, status, accrual, user_id, created_at, updated_at
	from loyalty.order
	where user_id = $1
	order by created_at;
`

// ListOrders returns list of orders for a given user.
func (r *Repo) ListOrders(ctx context.Context, userID int64) ([]*model.Order, error) {
	var orders []*model.Order
	err := r.db.QuerySlice(ctx, &orders, listOrders, userID)
	return orders, errs.Wrap(err, "query slice")
}

const findOrderByNumber = `
	select number, status, accrual, user_id, created_at, updated_at
	from loyalty.order
	where number = $1;
`

// FindOrderByNumber finds an order by its number.
func (r *Repo) FindOrderByNumber(ctx context.Context, orderNumber string) (*model.Order, error) {
	var order model.Order
	err := r.db.QueryRow(ctx, &order, findOrderByNumber, orderNumber)
	return &order, errs.Wrap(err, "query struct")
}

const findPollableOrders = `
	select order_number as number, attempt, last_polled, next_poll
	from loyalty.order o
	join loyalty.poll_accrual pa
		on o.number = pa.order_number
	where pa.next_poll <= current_timestamp
		and o.status in ('NEW', 'PROCESSING');
`

// FindPollableOrders finds orders that may be polled for updates.
func (r *Repo) FindPollableOrders(ctx context.Context) ([]*model.PollableOrder, error) {
	var orders []*model.PollableOrder
	err := r.db.QuerySlice(ctx, &orders, findPollableOrders)
	return orders, errs.Wrap(err, "query slice")
}
