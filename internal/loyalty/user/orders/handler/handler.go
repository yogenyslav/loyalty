// Package handler routes HTTP requests for orders.
package handler

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"
)

type orderController interface {
	ProcessOrder(ctx context.Context, orderNumber string, userID int64) (string, error)
	ListOrders(ctx context.Context, userID int64) ([]model.OrderDto, error)
}

// Handler provides methods to handle HTTP requests for orders.
type Handler struct {
	oc orderController
}

// New creates a new Handler instance.
func New(oc orderController) *Handler {
	return &Handler{oc: oc}
}
