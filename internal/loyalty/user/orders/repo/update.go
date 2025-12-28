package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

const updateOrderStatus = `
	update loyalty.order
	set status = $1, 
		updated_at = current_timestamp
	where number = $2;
`

// UpdateOrderStatus updates the status of an existing order.
func (r *Repo) UpdateOrderStatus(ctx context.Context, orderNumber string, status model.Status) error {
	rowsAffected, err := r.db.Exec(ctx, updateOrderStatus, status, orderNumber)
	if err != nil {
		return errs.Wrap(err, "exec query")
	}

	if rowsAffected == 0 {
		return errs.Wrap(pgx.ErrNoRows)
	}

	return nil
}

const updatePollingSchedule = `
	update loyalty.poll_accrual
	set last_polled = current_timestamp,
		attempt = attempt + 1,
		next_poll = current_timestamp + interval '1 second' * $1
	where order_number = $2;
`

// UpdatePollingSchedule updates the polling schedule for a given order.
func (r *Repo) UpdatePollingSchedule(ctx context.Context, orderNumber string, backoffSeconds int) error {
	_, err := r.db.Exec(ctx, updatePollingSchedule, backoffSeconds, orderNumber)
	return errs.Wrap(err, "exec query")
}
