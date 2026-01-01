package controller

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/loyalty/internal/accrual"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"
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

		ctrl := New(orderRepo, balanceUpdater, accrualService)
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

		orderRepo.EXPECT().
			BeginTx(ctx, gomock.Any()).
			Return(ctx, nil)

		balanceUpdater.EXPECT().
			UpdateBalanceAccrual(ctx, userID, 123.45).
			Return(nil)

		orderRepo.EXPECT().
			UpdateOrderAfterPolling(ctx, gomock.Any()).
			Return(nil)

		orderRepo.EXPECT().
			CommitTx(ctx).
			Return(nil)

		orderRepo.EXPECT().
			RollbackTx(ctx)

		orderNumber, err := ctrl.ProcessOrder(ctx, orderNumber, userID)
		require.NoError(t, err)
		require.Equal(t, "2844830162", orderNumber)
	})

	t.Run("SQL error on inserting new order", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService)
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

		ctrl := New(orderRepo, balanceUpdater, accrualService)
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

		ctrl := New(orderRepo, balanceUpdater, accrualService)
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

		ctrl := New(orderRepo, balanceUpdater, accrualService)
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

		orderRepo.EXPECT().
			UpdateOrderAfterPolling(ctx, &accrual.OrderAccrual{
				Order:  orderNumber,
				Status: model.StatusProcessing,
			}).
			Return(nil)

		orderNumber, err := ctrl.ProcessOrder(ctx, orderNumber, userID)
		require.NoError(t, err)
		require.Equal(t, "2844830162", orderNumber)
	})

	t.Run("Order is created first time, not found in accrual system", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService)
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

		ctrl := New(orderRepo, balanceUpdater, accrualService)
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

		ctrl := New(orderRepo, balanceUpdater, accrualService)
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

		ctrl := New(orderRepo, balanceUpdater, accrualService)
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

		ctrl := New(orderRepo, balanceUpdater, accrualService)
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

	t.Run("Pollable orders are processed", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService)
		ctx := t.Context()

		orderNumber := "2844830162"

		today := time.Now()
		orders := []*model.PollableOrder{
			{
				Number:     orderNumber,
				LastPolled: today,
				NextPoll:   today,
			},
		}

		accrualService.EXPECT().
			PollInterval().
			Return(1)

		orderRepo.EXPECT().
			FindPollableOrders(ctx).
			Return(orders, nil)

		accrualService.EXPECT().
			Poll(ctx, gomock.Any(), gomock.Any())

		errCh := make(chan error)
		go ctrl.PollOrders(ctx, errCh)

		go func() {
			for err := range errCh {
				require.NoError(t, err)
			}
		}()
		time.Sleep(time.Second * 2)
	})
}

func TestController_processPollResult(t *testing.T) {
	t.Parallel()

	t.Run("Process poll results with error", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService)
		ctx := t.Context()

		orderNumber := "2844830162"

		accrualService.EXPECT().
			PollInterval().
			Return(1)

		orderRepo.EXPECT().
			UpdatePollingSchedule(ctx, orderNumber, gomock.Any()).
			Return(nil)

		errCh := make(chan error)
		go ctrl.processPollResult(ctx, orderNumber, 0, accrual.PollResult{
			OrderAccrual: nil,
			Err:         errors.New("polling error"),
		}, errCh)

		go func() {
			for err := range errCh {
				require.Error(t, err)
			}
		}()
		time.Sleep(time.Second * 2)
	})

	t.Run("Process poll results with accrual info", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService)
		ctx := t.Context()

		orderNumber := "2844830162"

		orderRepo.EXPECT().
			BeginTx(ctx, gomock.Any()).
			Return(ctx, nil)

		orderRepo.EXPECT().
			UpdateOrderAfterPolling(ctx, gomock.Any()).
			Return(nil)

		orderRepo.EXPECT().
			DeletePollingSchedule(ctx, orderNumber).
			Return(nil)

		orderRepo.EXPECT().
			FindOrderByNumber(ctx, orderNumber).
			Return(&model.Order{
				Number: orderNumber,
				UserID: int64(1),
			}, nil)

		balanceUpdater.EXPECT().
			UpdateBalanceAccrual(ctx, int64(1), 123.45).
			Return(nil)

		orderRepo.EXPECT().
			CommitTx(ctx).
			Return(nil)

		orderRepo.EXPECT().
			RollbackTx(ctx)

		errCh := make(chan error)
		go ctrl.processPollResult(ctx, orderNumber, 0, accrual.PollResult{
			OrderAccrual: &accrual.OrderAccrual{
				Order:   orderNumber,
				Status:  model.StatusProcessed,
				Accrual: 123.45,
			},
			Err: nil,
		}, errCh)

		go func() {
			for err := range errCh {
				require.NoError(t, err)
			}
		}()
		time.Sleep(time.Second * 2)
	})

	t.Run("Process poll results with nil accrual info", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService)
		ctx := t.Context()

		orderNumber := "2844830162"

		accrualService.EXPECT().
			PollInterval().
			Return(1)

		orderRepo.EXPECT().
			UpdatePollingSchedule(ctx, orderNumber, gomock.Any()).
			Return(nil)

		errCh := make(chan error)
		go ctrl.processPollResult(ctx, orderNumber, 0, accrual.PollResult{
			OrderAccrual: nil,
			Err:         nil,
		}, errCh)

		go func() {
			for err := range errCh {
				require.NoError(t, err)
			}
		}()
		time.Sleep(time.Second * 2)
	})

	t.Run("Process poll results with non-final status", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService)
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

		errCh := make(chan error)
		go ctrl.processPollResult(ctx, orderNumber, 0, accrual.PollResult{
			OrderAccrual: &accrual.OrderAccrual{
				Order:  orderNumber,
				Status: model.StatusProcessing,
			},
			Err: nil,
		}, errCh)

		go func() {
			for err := range errCh {
				require.NoError(t, err)
			}
		}()
		time.Sleep(time.Second * 2)
	})

	t.Run("Process poll results with error on updating schedule", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService)
		ctx := t.Context()

		orderNumber := "2844830162"

		accrualService.EXPECT().
			PollInterval().
			Return(1)

		orderRepo.EXPECT().
			UpdatePollingSchedule(ctx, orderNumber, gomock.Any()).
			Return(errors.New("error updating schedule"))

		errCh := make(chan error)
		go ctrl.processPollResult(ctx, orderNumber, 0, accrual.PollResult{
			OrderAccrual: nil,
			Err:         errors.New("polling error"),
		}, errCh)

		go func() {
			for err := range errCh {
				require.Error(t, err)
			}
		}()
		time.Sleep(time.Second * 2)
	})

	t.Run("Process poll results with error on processing with data", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService)
		ctx := t.Context()

		orderNumber := "2844830162"

		orderRepo.EXPECT().
			BeginTx(ctx, gomock.Any()).
			Return(ctx, nil)

		orderRepo.EXPECT().
			UpdateOrderAfterPolling(ctx, gomock.Any()).
			Return(nil)

		orderRepo.EXPECT().
			FindOrderByNumber(ctx, orderNumber).
			Return(&model.Order{
				Number: orderNumber,
				UserID: int64(1),
			}, nil)

		balanceUpdater.EXPECT().
			UpdateBalanceAccrual(ctx, int64(1), 123.45).
			Return(errors.New("error updating balance"))

		orderRepo.EXPECT().
			DeletePollingSchedule(ctx, orderNumber).
			Return(nil)

		orderRepo.EXPECT().
			RollbackTx(ctx)

		errCh := make(chan error)
		go ctrl.processPollResult(ctx, orderNumber, 0, accrual.PollResult{
			OrderAccrual: &accrual.OrderAccrual{
				Order:   orderNumber,
				Status:  model.StatusProcessed,
				Accrual: 123.45,
			},
			Err: nil,
		}, errCh)

		go func() {
			for err := range errCh {
				require.Error(t, err)
			}
		}()
		time.Sleep(time.Second * 2)
	})

	t.Run("Process poll results with error on updating order after polling", func(t *testing.T) {
		t.Parallel()

		orderRepo := mocks.NewMockorderRepo(gomock.NewController(t))
		balanceUpdater := mocks.NewMockbalanceUpdater(gomock.NewController(t))
		accrualService := mocks.NewMockAccrualService(gomock.NewController(t))

		ctrl := New(orderRepo, balanceUpdater, accrualService)
		ctx := t.Context()

		orderNumber := "2844830162"

		orderRepo.EXPECT().
			BeginTx(ctx, gomock.Any()).
			Return(ctx, nil)

		orderRepo.EXPECT().
			UpdateOrderAfterPolling(ctx, gomock.Any()).
			Return(errors.New("error updating order after polling"))

		orderRepo.EXPECT().
			RollbackTx(ctx)

		errCh := make(chan error)
		go ctrl.processPollResult(ctx, orderNumber, 0, accrual.PollResult{
			OrderAccrual: &accrual.OrderAccrual{
				Order:   orderNumber,
				Status:  model.StatusProcessed,
				Accrual: 123.45,
			},
			Err: nil,
		}, errCh)

		go func() {
			for err := range errCh {
				require.Error(t, err)
			}
		}()
		time.Sleep(time.Second * 2)
	})
}
