package controller

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
	"github.com/yogenyslav/loyalty/pkg/luhn"
)

// ProcessOrder processes an order for a given user.
// Returns orderNumber if it is created firt time or empty string if it already exists.
// In case of error, returns non-nil error.
func (ctrl *Controller) ProcessOrder(ctx context.Context, orderNumber string, userID int64) (string, error) {
	if !luhn.Validate(orderNumber) {
		return "", errs.Wrap(luhn.ErrInvalidNumber, "validate order number")
	}

	order := &model.Order{
		Number: orderNumber,
		UserID: userID,
	}
	err := ctrl.or.InsertOrder(ctx, order)
	if err == nil { // no error means order is created first time
		err = ctrl.processNewOrder(ctx, userID, orderNumber)
		if err != nil {
			return "", errs.Wrap(err, "process new order")
		}
		return orderNumber, nil
	}
	// some error occurred

	// if not unique violation, its fatal insert error
	if !errs.CheckUniqueViloation(err) {
		return "", errs.Wrap(err, "insert order")
	}

	// order already exists, check the creator
	order, err = ctrl.or.FindOrderByNumber(ctx, orderNumber)
	if err != nil {
		return "", errs.Wrap(err, "find order by number")
	}

	if order.UserID != userID {
		return "", errs.Wrap(errs.ErrOrderCreatedByAnotherUser, "order belongs to another user")
	}

	return "", nil
}

func (ctrl *Controller) processNewOrder(ctx context.Context, userID int64, orderNumber string) error {
	accrualInfo, err := ctrl.accrualService.ProcessOrder(ctx, orderNumber)
	if err != nil {
		return errs.Wrap(err, "process order in accrual service")
	}

	if accrualInfo == nil || accrualInfo.Status == model.StatusProcessing ||
		accrualInfo.Status == model.StatusRegistered {
		err = ctrl.or.SchedulePolling(ctx, orderNumber)
		if err != nil {
			return errs.Wrap(err, "schedule order polling")
		}
	}

	// order is not in the accrual system
	if accrualInfo == nil {
		return nil
	}

	// update according to the real balance and status
	var updatedBalance bool
	if accrualInfo.Accrual > 0 && accrualInfo.Status == model.StatusProcessed {
		tx, err := ctrl.or.BeginTx(ctx)
		if err != nil {
			return errs.Wrap(err, "begin tx for new order processing")
		}
		defer ctrl.or.RollbackTx(ctx) //nolint:errcheck // nothing we can do

		ctx = tx

		err = ctrl.br.UpdateBalanceAccrual(ctx, userID, accrualInfo.Accrual)
		if err != nil {
			return errs.Wrap(err, "update balance accrual")
		}
		updatedBalance = true
	}

	err = ctrl.or.UpdateOrderAfterPolling(ctx, accrualInfo)
	if err != nil {
		return errs.Wrap(err, "update order status")
	}

	if updatedBalance {
		err = ctrl.or.CommitTx(ctx)
		if err != nil {
			return errs.Wrap(err, "commit tx for new order processing")
		}
	}

	return nil
}
