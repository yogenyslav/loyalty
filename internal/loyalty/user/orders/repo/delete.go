package repo

import (
	"context"

	"github.com/yogenyslav/loyalty/pkg/errs"
)

const deletePollingSchedule = `
	delete from loyalty.poll_accrual
	where order_number = $1;
`

// DeletePollingSchedule deletes the polling schedule for a given order.
func (r *Repo) DeletePollingSchedule(ctx context.Context, orderNumber string) error {
	_, err := r.db.Exec(ctx, deletePollingSchedule, orderNumber)
	return errs.Wrap(err, "exec query")
}
