package controller

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

// ListOrders returns list of orders for a given user.
func (ctrl *Controller) ListOrders(ctx context.Context, userID int64) ([]model.OrderDto, error) {
	orders, err := ctrl.or.ListOrders(ctx, userID)
	if err != nil {
		return nil, errs.Wrap(err, "list orders")
	}

	var ordersDto []model.OrderDto
	for _, order := range orders {
		ordersDto = append(ordersDto, *order.ToDto())
	}

	return ordersDto, nil
}
