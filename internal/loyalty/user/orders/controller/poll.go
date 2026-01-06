package controller

import (
	"context"
	"log/slog"
	"time"

	"github.com/yogenyslav/loyalty/internal/accrual"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"
	"github.com/yogenyslav/loyalty/pkg/database"
	"github.com/yogenyslav/loyalty/pkg/errs"
	"golang.org/x/sync/errgroup"
)

const maxBackoffSeconds = 600

// PollOrders polls the accrual system for order updates.
func (ctrl *Controller) PollOrders(ctx context.Context, errCh chan<- error) {
	timer := time.NewTicker(time.Second * time.Duration(ctrl.accrualService.PollInterval()))
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			slog.Info("polling orders")

			pollableOrders, err := ctrl.or.FindPollableOrders(ctx)
			if err != nil {
				errCh <- errs.Wrap(err, "find pollable orders")
			}

			orders := make(chan string, len(pollableOrders))
			results := make(chan *accrual.OrderAccrual, len(pollableOrders))

			err = ctrl.pollOnce(ctx, orders, results, pollableOrders)
			errCh <- errs.Wrap(err, "poll once")
		}
	}
}

func (ctrl *Controller) pollOnce(
	ctx context.Context,
	orders chan string,
	results chan *accrual.OrderAccrual,
	pollableOrders map[string]*model.PollableOrder,
) error {
	if len(pollableOrders) == 0 {
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)
	numWorkers := ctrl.accrualService.NumWorkers()

	// explicitly set limit to numWorkers, although goroutines are only created in range over numWorkers
	g.SetLimit(numWorkers)

	for range numWorkers {
		g.Go(func() error {
			return ctrl.accrualService.Poll(ctx, orders, results)
		})
	}

	for _, order := range pollableOrders {
		orders <- order.Number
	}
	close(orders)

	// store pollingError to process the case when polling failed
	var pollingError error
	if pollingError = g.Wait(); pollingError != nil {
		slog.Error("error polling orders", slog.Any("error", errs.Wrap(pollingError)))
	}

processResults:
	for {
		select {
		// all orders that have been polled and processed successfully have to be marked as 'no pollingError' so that they are not rescheduled
		case result := <-results:
			pollAttempt := pollableOrders[result.Order].Attempt
			g.Go(func() error {
				return ctrl.processPollResult(
					ctx,
					result.Order,
					pollAttempt,
					result,
					nil,
				)
			})
			// also delete them from map to process the rest that are meant to be failed
			delete(pollableOrders, result.Order)
		default:
			// processing failed (or those that were after the first failed)
			if pollingError != nil {
				for _, order := range pollableOrders {
					g.Go(func() error {
						return ctrl.processPollResult(
							ctx,
							order.Number,
							order.Attempt,
							nil,
							pollingError,
						)
					})
				}
			}
			break processResults
		}
	}

	return errs.Wrap(g.Wait(), "process poll results")
}

func (ctrl *Controller) processPollResult(
	ctx context.Context,
	orderNumber string,
	attempt int,
	result *accrual.OrderAccrual,
	pollingError error,
) error {
	// have polling error -> reschedule with backoff
	if pollingError != nil {
		backoffSeconds := getBackoffSeconds(ctrl.accrualService.PollInterval(), attempt)
		err := ctrl.or.UpdatePollingSchedule(ctx, orderNumber, backoffSeconds)
		return errs.Wrap(err, "update polling schedule after polling error")
	}

	// only if accrual info is returned
	if result != nil {
		err := ctrl.uow.WithTx(ctx, database.TxLevelSerializable, func(ctx context.Context) error {
			return ctrl.processWithData(ctx, orderNumber, result)
		})
		return errs.Wrap(err, "process poll result with data")
	}

	// no error, no data -> reschedule with base backoff
	backoffSeconds := ctrl.accrualService.PollInterval()
	err := ctrl.or.UpdatePollingSchedule(ctx, orderNumber, backoffSeconds)
	return errs.Wrap(err, "update polling schedule after polling")
}

func (ctrl *Controller) processWithData(
	ctx context.Context,
	orderNumber string,
	accrualInfo *accrual.OrderAccrual,
) error {
	if accrualInfo == nil {
		return nil
	}

	// update order status and accrual
	err := ctrl.or.UpdateOrderAfterPolling(ctx, accrualInfo)
	if err != nil {
		return errs.Wrap(err, "update order after polling")
	}

	switch accrualInfo.Status {
	// polling returned final status, delete from polling schedule
	case model.StatusProcessed, model.StatusInvalid:
		err = ctrl.or.DeletePollingSchedule(ctx, orderNumber)
		if err != nil {
			return errs.Wrap(err, "delete polling schedule after final status")
		}

		if accrualInfo.Status == model.StatusInvalid {
			return nil
		}

		order, err := ctrl.or.FindOrderByNumber(ctx, orderNumber)
		if err != nil {
			return errs.Wrap(err, "find order by number after polling")
		}

		err = ctrl.br.UpdateBalanceAccrual(ctx, order.UserID, accrualInfo.Accrual)
		if err != nil {
			return errs.Wrap(err, "update balance accrual after polling")
		}

		return nil
	default:
		backoffSeconds := ctrl.accrualService.PollInterval()
		err = ctrl.or.UpdatePollingSchedule(ctx, orderNumber, backoffSeconds)
		return errs.Wrap(err, "reschedule polling for non-final status")
	}
}
