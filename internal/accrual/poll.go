package accrual

import (
	"context"

	"github.com/yogenyslav/loyalty/pkg/errs"
)

// PollResult represents the result of polling an order's accrual information.
type PollResult struct {
	*OrderAccrual
	Err error
}

// Poll polls accrual service for incoming orders and sends the results to channel.
func (ac *Service) Poll(ctx context.Context, orders <-chan string, result chan<- PollResult) {
	for range ac.cfg.NumWorkers {
		go func() {
			for orderID := range orders {
				accrualInfo, err := ac.ProcessOrder(ctx, orderID)
				result <- PollResult{OrderAccrual: accrualInfo, Err: errs.Wrap(err, "process order while polling")}
			}
		}()
	}
}
