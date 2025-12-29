package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/loyalty/internal/accrual"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

const updateOrderAfterPolling = `
	update loyalty.order
	set status = $2, 
		accrual = $3,
		updated_at = current_timestamp
	where number = $1;
`

// UpdateOrderAfterPolling updates the order with data received from the accrual service.
func (r *Repo) UpdateOrderAfterPolling(ctx context.Context, data *accrual.AccrualInfo) error {
	if data == nil {
		return nil
	}

	rowsAffected, err := r.db.Exec(ctx, updateOrderAfterPolling, data.Order, data.Status, data.Accrual)
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
