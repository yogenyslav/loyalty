// Package controller implements logic layer for order management.
package controller

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/accrual"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"
	"github.com/yogenyslav/loyalty/pkg/database"
)

//go:generate mockgen -destination=../../../../../tests/mocks/order_repo.go -package=mocks . orderRepo,balanceUpdater
type orderRepo interface {
	InsertOrder(ctx context.Context, order *model.Order) error
	ListOrders(ctx context.Context, userID int64) ([]*model.Order, error)
	FindOrderByNumber(ctx context.Context, orderNumber string) (*model.Order, error)
	FindPollableOrders(ctx context.Context) ([]*model.PollableOrder, error)
	UpdateOrderAfterPolling(ctx context.Context, data *accrual.AccrualInfo) error
	SchedulePolling(ctx context.Context, orderNumber string) error
	UpdatePollingSchedule(ctx context.Context, orderNumber string, backoffSeconds int) error
	DeletePollingSchedule(ctx context.Context, orderNumber string) error

	BeginTx(ctx context.Context, level database.TxLevel) (context.Context, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error
}

type balanceUpdater interface {
	UpdateBalanceAccrual(ctx context.Context, userID int64, accrual float64) error
}

// AccrualService defines the interface for interacting with the accrual system.
//
//go:generate mockgen -destination=../../../../../tests/mocks/accrual_service.go -package=mocks . AccrualService
type AccrualService interface {
	ProcessOrder(ctx context.Context, orderNumber string) (*accrual.AccrualInfo, error)
	Poll(ctx context.Context, orders <-chan string, result chan<- accrual.PollResult)
	PollInterval() int
}

// Controller provides methods for order management logic.
type Controller struct {
	or             orderRepo
	br             balanceUpdater
	accrualService AccrualService
}

// New creates a new Controller instance.
func New(or orderRepo, br balanceUpdater, accrual AccrualService) *Controller {
	return &Controller{
		or:             or,
		br:             br,
		accrualService: accrual,
	}
}
