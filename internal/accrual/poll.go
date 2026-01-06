package accrual

import (
	"context"

	"github.com/yogenyslav/loyalty/pkg/errs"
)

// Poll polls accrual service for incoming orders and sends the results to channel.
func (ac *Service) Poll(ctx context.Context, orders <-chan string, result chan<- *OrderAccrual) error {
	for orderID := range orders {
		accrualInfo, err := ac.ProcessOrder(ctx, orderID)
		if err != nil {
			return errs.Wrap(err, "process order in accrual system")
		}
		result <- accrualInfo
	}
	return nil
}
