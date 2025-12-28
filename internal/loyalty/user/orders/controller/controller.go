// Package controller implements logic layer for order management.
package controller

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/accrual"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"
)

type orderRepo interface {
	InsertOrder(ctx context.Context, order *model.Order) error
	ListOrders(ctx context.Context, userID int64) ([]*model.Order, error)
	FindOrderByNumber(ctx context.Context, orderNumber string) (*model.Order, error)
	FindPollableOrders(ctx context.Context) ([]*model.PollableOrder, error)
	SchedulePolling(ctx context.Context, orderNumber string) error
	UpdateOrderStatus(ctx context.Context, orderNumber string, status model.Status) error
	UpdatePollingSchedule(ctx context.Context, orderNumber string, backoffSeconds int) error
	DeletePollingSchedule(ctx context.Context, orderNumber string) error
}

type balanceRepo interface {
	UpdateBalanceAccrual(ctx context.Context, userID int64, accrual float64) error
}

// AccrualService defines the interface for interacting with the accrual system.
type AccrualService interface {
	ProcessOrder(ctx context.Context, orderNumber string) (*accrual.AccrualInfo, error)
	Poll(ctx context.Context, orders <-chan string, result chan<- accrual.PollResult)
	PollInterval() int
}

// Controller provides methods for order management logic.
type Controller struct {
	or             orderRepo
	br             balanceRepo
	accrualService AccrualService
}

// New creates a new Controller instance.
func New(or orderRepo, br balanceRepo, accrual AccrualService) *Controller {
	return &Controller{
		or:             or,
		br:             br,
		accrualService: accrual,
	}
}
