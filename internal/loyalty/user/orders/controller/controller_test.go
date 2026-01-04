package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/loyalty/internal/accrual"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"
	"github.com/yogenyslav/loyalty/pkg/database"
	"github.com/yogenyslav/loyalty/pkg/errs"
	"github.com/yogenyslav/loyalty/pkg/luhn"
	"github.com/yogenyslav/loyalty/tests/mocks"
	"go.uber.org/mock/gomock"
)

func TestController_Process(t *testing.T) {
	t.Parallel()

	t.Run("New order is processed for the first time", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()
		orderNumber := "2844830162"
		userID := int64(1)

		orderRepo.EXPECT().
			InsertOrder(ctx, gomock.Any()).
			Return(nil)

		accrualService.EXPECT().
			ProcessOrder(ctx, orderNumber).
			Return(&accrual.OrderAccrual{
				Order:   orderNumber,
				Status:  model.StatusProcessed,
				Accrual: 123.45,
			}, nil)

		uow.EXPECT().
			WithTx(ctx, database.TxLevelSerializable, gomock.Any()).
			DoAndReturn(func(ctx context.Context, level database.TxLevel, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})

		orderRepo.EXPECT().
			UpdateOrderAfterPolling(ctx, gomock.Any()).
			Return(nil)

		balanceUpdater.EXPECT().
			UpdateBalanceAccrual(ctx, userID, 123.45).
			Return(nil)

		orderNumber, err := ctrl.ProcessOrder(ctx, orderNumber, userID)
		require.NoError(t, err)
		require.Equal(t, "2844830162", orderNumber)
	})

	t.Run("SQL error on inserting new order", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()
		orderNumber := "2844830162"
		userID := int64(1)

		orderRepo.EXPECT().
			InsertOrder(ctx, gomock.Any()).
			Return(errors.New("sql error"))

		orderNumber, err := ctrl.ProcessOrder(ctx, orderNumber, userID)
		require.Error(t, err)
		require.Equal(t, "", orderNumber)
	})

	t.Run("Accrual service error on processing new order", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()
		orderNumber := "2844830162"
		userID := int64(1)

		orderRepo.EXPECT().
			InsertOrder(ctx, gomock.Any()).
			Return(nil)

		accrualService.EXPECT().
			ProcessOrder(ctx, orderNumber).
			Return(nil, errors.New("accrual service error"))

		_, err := ctrl.ProcessOrder(ctx, orderNumber, userID)
		require.Error(t, err)
	})

	t.Run("Invalid order number", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()
		orderNumber := "12345"

		orderNumber, err := ctrl.ProcessOrder(ctx, orderNumber, 1)
		require.ErrorIs(t, err, luhn.ErrInvalidNumber)
		require.Equal(t, "", orderNumber)
	})

	t.Run("Order is created first time, accrual is not calculated yet", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()
		orderNumber := "2844830162"
		userID := int64(1)

		orderRepo.EXPECT().
			InsertOrder(ctx, gomock.Any()).
			Return(nil)

		accrualService.EXPECT().
			ProcessOrder(ctx, orderNumber).
			Return(&accrual.OrderAccrual{
				Order:  orderNumber,
				Status: model.StatusProcessing,
			}, nil)

		orderRepo.EXPECT().
			SchedulePolling(ctx, orderNumber).
			Return(nil)

		uow.EXPECT().
			WithTx(ctx, database.TxLevelSerializable, gomock.Any())

		orderNumber, err := ctrl.ProcessOrder(ctx, orderNumber, userID)
		require.NoError(t, err)
		require.Equal(t, "2844830162", orderNumber)
	})

	t.Run("Order is created first time, not found in accrual system", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()
		orderNumber := "2844830162"
		userID := int64(1)

		orderRepo.EXPECT().
			InsertOrder(ctx, gomock.Any()).
			Return(nil)

		accrualService.EXPECT().
			ProcessOrder(ctx, orderNumber).
			Return(nil, nil)

		orderRepo.EXPECT().
			SchedulePolling(ctx, orderNumber).
			Return(nil)

		orderNumber, err := ctrl.ProcessOrder(ctx, orderNumber, userID)
		require.NoError(t, err)
		require.Equal(t, "2844830162", orderNumber)
	})

	t.Run("Order already exists for the same user", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()
		orderNumber := "2844830162"
		userID := int64(1)

		orderRepo.EXPECT().
			InsertOrder(ctx, gomock.Any()).
			Return(&pgconn.PgError{Code: pgerrcode.UniqueViolation})

		orderRepo.EXPECT().
			FindOrderByNumber(ctx, orderNumber).
			Return(&model.Order{
				Number: orderNumber,
				UserID: userID,
			}, nil)

		orderNumber, err := ctrl.ProcessOrder(ctx, orderNumber, userID)
		require.NoError(t, err)
		require.Equal(t, "", orderNumber)
	})

	t.Run("Order is already created by another user", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()
		orderNumber := "2844830162"
		userID := int64(1)

		orderRepo.EXPECT().
			InsertOrder(ctx, gomock.Any()).
			Return(&pgconn.PgError{Code: pgerrcode.UniqueViolation})

		orderRepo.EXPECT().
			FindOrderByNumber(ctx, orderNumber).
			Return(&model.Order{
				Number: orderNumber,
				UserID: 2,
			}, nil)

		orderNumber, err := ctrl.ProcessOrder(ctx, orderNumber, userID)
		require.ErrorIs(t, err, errs.ErrOrderCreatedByAnotherUser)
		require.Equal(t, "", orderNumber)
	})
}

func TestController_ListOrders(t *testing.T) {
	t.Parallel()

	t.Run("Orders are listed", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()
		userID := int64(1)
		orders := []*model.Order{
			{
				Number: "2844830162",
				UserID: userID,
			},
			{
				Number: "28448301623",
				UserID: userID,
			},
		}

		orderRepo.EXPECT().
			ListOrders(ctx, userID).
			Return(orders, nil)

		resp, err := ctrl.ListOrders(ctx, userID)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp, 2)
	})

	t.Run("Error from repo", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()
		userID := int64(1)

		orderRepo.EXPECT().
			ListOrders(ctx, userID).
			Return(nil, errors.New("error"))

		resp, err := ctrl.ListOrders(ctx, userID)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestController_PollOrders(t *testing.T) {
	t.Parallel()

	t.Run("Smoke test for polling orders", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()

		accrualService.EXPECT().
			PollInterval().
			Return(1)

		orderRepo.EXPECT().
			FindPollableOrders(ctx).
			Return(map[string]*model.PollableOrder{}, nil)

		errCh := make(chan error, 1)
		go ctrl.PollOrders(ctx, errCh)

		<-errCh
	})
}

func TestController_pollOnce(t *testing.T) {
	t.Parallel()

	t.Run("Poll once with no pollable orders", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()

		err := ctrl.pollOnce(ctx, nil, nil, map[string]*model.PollableOrder{})
		require.NoError(t, err)
	})

	t.Run("Poll once with pollable orders", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()

		pollableOrders := map[string]*model.PollableOrder{
			"2844830162": {
				Number:  "2844830162",
				Attempt: 0,
			},
			"28448301623": {
				Number:  "28448301623",
				Attempt: 1,
			},
		}
		orders := make(chan string, len(pollableOrders))
		results := make(chan *accrual.OrderAccrual, len(pollableOrders))

		accrualService.EXPECT().
			NumWorkers().
			Return(1)

		accrualService.EXPECT().
			Poll(gomock.Any(), orders, results).
			DoAndReturn(func(ctx context.Context, orders <-chan string, results chan<- *accrual.OrderAccrual) error {
				for order := range orders {
					results <- &accrual.OrderAccrual{
						Order:  order,
						Status: model.StatusProcessing,
					}
				}
				return nil
			})

		uow.EXPECT().
			WithTx(gomock.Any(), database.TxLevelSerializable, gomock.Any()).
			Return(nil).
			Times(len(pollableOrders))

		err := ctrl.pollOnce(ctx, orders, results, pollableOrders)
		require.NoError(t, err)
	})
}

func TestController_processPollResult(t *testing.T) {
	t.Parallel()

	t.Run("Process poll results with error", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()

		orderNumber := "2844830162"

		accrualService.EXPECT().
			PollInterval().
			Return(1)

		orderRepo.EXPECT().
			UpdatePollingSchedule(ctx, orderNumber, gomock.Any()).
			Return(nil)

		err := ctrl.processPollResult(ctx, orderNumber, 0, nil, errors.New("polling error"))
		require.NoError(t, err)
	})

	t.Run("Process poll results with accrual info", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()

		orderNumber := "2844830162"

		uow.EXPECT().
			WithTx(ctx, database.TxLevelSerializable, gomock.Any()).
			DoAndReturn(func(ctx context.Context, level database.TxLevel, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})

		orderRepo.EXPECT().
			UpdateOrderAfterPolling(ctx, gomock.Any()).
			Return(nil)

		orderRepo.EXPECT().
			DeletePollingSchedule(ctx, gomock.Any()).
			Return(nil)

		orderRepo.EXPECT().
			FindOrderByNumber(ctx, orderNumber).
			Return(&model.Order{
				Number: orderNumber,
				UserID: 1,
			}, nil)

		balanceUpdater.EXPECT().
			UpdateBalanceAccrual(ctx, int64(1), 123.45).
			Return(nil)

		err := ctrl.processPollResult(ctx, orderNumber, 0, &accrual.OrderAccrual{
			Order:   orderNumber,
			Status:  model.StatusProcessed,
			Accrual: 123.45,
		}, nil)
		require.NoError(t, err)
	})

	t.Run("Process poll results with nil accrual info", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()

		orderNumber := "2844830162"

		accrualService.EXPECT().
			PollInterval().
			Return(1)

		orderRepo.EXPECT().
			UpdatePollingSchedule(ctx, orderNumber, gomock.Any()).
			Return(nil)

		err := ctrl.processPollResult(ctx, orderNumber, 0, nil, nil)
		require.NoError(t, err)
	})

	t.Run("Process poll results with non-final status", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()

		orderNumber := "2844830162"

		accrualService.EXPECT().
			PollInterval().
			Return(1)

		orderRepo.EXPECT().
			UpdateOrderAfterPolling(ctx, gomock.Any()).
			Return(nil)

		orderRepo.EXPECT().
			UpdatePollingSchedule(ctx, orderNumber, gomock.Any()).
			Return(nil)

		uow.EXPECT().
			WithTx(ctx, database.TxLevelSerializable, gomock.Any()).
			DoAndReturn(func(ctx context.Context, level database.TxLevel, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})

		err := ctrl.processPollResult(ctx, orderNumber, 0, &accrual.OrderAccrual{
			Order:  orderNumber,
			Status: model.StatusProcessing,
		}, nil)
		require.NoError(t, err)
	})

	t.Run("Process poll results with error on updating schedule", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()

		orderNumber := "2844830162"

		accrualService.EXPECT().
			PollInterval().
			Return(1)

		orderRepo.EXPECT().
			UpdatePollingSchedule(ctx, orderNumber, gomock.Any()).
			Return(errors.New("error updating schedule"))

		orderRepo.EXPECT().
			UpdateOrderAfterPolling(ctx, gomock.Any()).
			Return(nil)

		uow.EXPECT().
			WithTx(ctx, database.TxLevelSerializable, gomock.Any()).
			DoAndReturn(func(ctx context.Context, level database.TxLevel, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})

		err := ctrl.processPollResult(ctx, orderNumber, 0, &accrual.OrderAccrual{
			Order:  orderNumber,
			Status: model.StatusProcessing,
		}, nil)
		require.Error(t, err)
	})

	t.Run("Process poll results with error on processing with data", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))
		uow := mocks.NewMockUnitOfWork(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService, uow)
		ctx := t.Context()

		orderNumber := "2844830162"

		uow.EXPECT().
			WithTx(ctx, database.TxLevelSerializable, gomock.Any()).
			Return(errors.New("error"))

		err := ctrl.processPollResult(ctx, orderNumber, 0, &accrual.OrderAccrual{
			Order:   orderNumber,
			Status:  model.StatusProcessed,
			Accrual: 123.45,
		}, nil)
		require.Error(t, err)
	})
}
