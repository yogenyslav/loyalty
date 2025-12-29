package controller

import (
	"context"
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
	switch {
	// polling error, increase backoff and update schedule
	case result.Err != nil:
		errCh <- errs.Wrap(result.Err, "poll order result")

		backoffSeconds := getBackoffSeconds(ctrl.accrualService.PollInterval(), attempt)
		err := ctrl.or.UpdatePollingSchedule(ctx, orderNumber, backoffSeconds)
		if err != nil {
			errCh <- errs.Wrap(err, "update polling schedule after polling error")
			return
		}

	// only if accrual info is returned
	case result.AccrualInfo != nil:
		err := ctrl.processWithData(ctx, orderNumber, result.AccrualInfo)
		if err != nil {
			errCh <- errs.Wrap(err, "process poll result with data")
			return
		}

	// for not final status, reset polling schedule to default interval
	case result.AccrualInfo == nil || result.AccrualInfo.Status == model.StatusProcessing || result.AccrualInfo.Status == model.StatusRegistered:
		backoffSeconds := ctrl.accrualService.PollInterval()
		err := ctrl.or.UpdatePollingSchedule(ctx, orderNumber, backoffSeconds)
		if err != nil {
			errCh <- errs.Wrap(err, "update polling schedule after polling")
			return
		}
	}
}

func (ctrl *Controller) processWithData(
	ctx context.Context,
	orderNumber string,
	accrualInfo *accrual.AccrualInfo,
) error {
	if accrualInfo == nil {
		return nil
	}

	if accrualInfo.Status == model.StatusProcessed || accrualInfo.Status == model.StatusInvalid {
		tx, err := ctrl.or.BeginTx(ctx)
		if err != nil {
			return errs.Wrap(err, "begin tx for polling processing")
		}
		defer ctrl.or.RollbackTx(ctx) //nolint:errcheck // nothing we can do

		// switch to tx context
		ctx = tx
	}

	// update order status and accrual
	err := ctrl.or.UpdateOrderAfterPolling(ctx, accrualInfo)
	if err != nil {
		return errs.Wrap(err, "update order after polling")
	}

	// polling returned final status, delete from polling schedule
	if accrualInfo.Status == model.StatusProcessed || accrualInfo.Status == model.StatusInvalid {
		err = ctrl.or.DeletePollingSchedule(ctx, orderNumber)
		if err != nil {
			return errs.Wrap(err, "delete polling schedule after final status")
		}

		order, err := ctrl.or.FindOrderByNumber(ctx, orderNumber)
		if err != nil {
			return errs.Wrap(err, "find order by number after polling")
		}

		err = ctrl.br.UpdateBalanceAccrual(ctx, order.UserID, accrualInfo.Accrual)
		if err != nil {
			return errs.Wrap(err, "update balance accrual after polling")
		}

		err = ctrl.or.CommitTx(ctx)
		if err != nil {
			return errs.Wrap(err, "commit tx for polling processing")
		}
	}

	return nil
}
