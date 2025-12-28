package controller

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/loyalty/internal/accrual"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
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
			log.Info().Msg("polling orders")

			pollableOrders, err := ctrl.or.FindPollableOrders(ctx)
			if err != nil {
				errCh <- errs.Wrap(err, "find pollable orders")
				continue
			}

			orders := make(chan string, len(pollableOrders))
			results := make(chan accrual.PollResult, len(pollableOrders))

			go ctrl.accrualService.Poll(ctx, orders, results)

			for _, order := range pollableOrders {
				orders <- order.Number
			}
			close(orders)

			for _, order := range pollableOrders {
				go ctrl.processPollResult(ctx, order.Number, order.Attempt, <-results, errCh)
			}
			close(results)
		}
	}
}

func (ctrl *Controller) processPollResult(
	ctx context.Context,
	orderNumber string,
	attempt int,
	result accrual.PollResult,
	errCh chan<- error,
) {
	// polling error, increase backoff and update schedule
	if result.Err != nil {
		errCh <- errs.Wrap(result.Err, "poll order result")

		backoffSeconds := getBackoffSeconds(ctrl.accrualService.PollInterval(), attempt)
		err := ctrl.or.UpdatePollingSchedule(ctx, orderNumber, backoffSeconds)
		if err != nil {
			errCh <- errs.Wrap(err, "update polling schedule after polling error")
		}
		return
	}

	// only if accrual info is returned
	if result.AccrualInfo != nil {
		// update order status
		// TODO: also update accrual amount.
		err := ctrl.or.UpdateOrderStatus(ctx, orderNumber, result.AccrualInfo.Status)
		if err != nil {
			errCh <- errs.Wrap(err, "update order status after polling")
			return
		}

		// polling returned final status, delete from polling schedule
		if result.AccrualInfo.Status == model.StatusProcessed || result.AccrualInfo.Status == model.StatusInvalid {
			err = ctrl.or.DeletePollingSchedule(ctx, orderNumber)
			if err != nil {
				errCh <- errs.Wrap(err, "delete polling schedule after final status")
				return
			}

			order, err := ctrl.or.FindOrderByNumber(ctx, orderNumber)
			if err != nil {
				errCh <- errs.Wrap(err, "find order by number after polling")
				return
			}

			err = ctrl.br.UpdateBalanceAccrual(ctx, order.UserID, result.AccrualInfo.Accrual)
			if err != nil {
				errCh <- errs.Wrap(err, "update balance accrual after polling")
			}
			return
		}
	}

	// for not final status, reset polling schedule
	if result.AccrualInfo == nil || result.AccrualInfo.Status == model.StatusProcessing {
		backoffSeconds := ctrl.accrualService.PollInterval()
		err := ctrl.or.UpdatePollingSchedule(ctx, orderNumber, backoffSeconds)
		if err != nil {
			errCh <- errs.Wrap(err, "update polling schedule after polling")
			return
		}
	}
}

func getBackoffSeconds(baseBackoff, attempt int) int {
	jitter := rand.IntN(baseBackoff)
	backoff := baseBackoff*(1<<attempt) + jitter
	if backoff > maxBackoffSeconds {
		return maxBackoffSeconds
	}
	return backoff
}
