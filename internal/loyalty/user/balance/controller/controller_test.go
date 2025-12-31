package controller

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/balance/model"
	"github.com/yogenyslav/loyalty/tests/mocks"
	"go.uber.org/mock/gomock"
)

func TestController_GetBalance(t *testing.T) {
	t.Parallel()

	t.Run("Balance is returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockbalanceRepo(gomock.NewController(t))

		ctrl := New(repo)

		userID := int64(1)
		balance := &model.Balance{
			Current:   150.0,
			Withdrawn: 50.0,
		}

		repo.EXPECT().
			FindBalanceByUserID(ctx, userID).
			Return(balance, nil)

		resp, err := ctrl.GetBalance(ctx, userID)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, *balance.ToDto(), *resp)
	})

	t.Run("Error from repo", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockbalanceRepo(gomock.NewController(t))

		ctrl := New(repo)
		userID := int64(1)

		repo.EXPECT().
			FindBalanceByUserID(ctx, userID).
			Return(nil, errors.New("error"))

		resp, err := ctrl.GetBalance(ctx, userID)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestController_ListWithdrawals(t *testing.T) {
	t.Parallel()

	t.Run("Withdrawals are returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockbalanceRepo(gomock.NewController(t))

		ctrl := New(repo)

		today := time.Now()
		userID := int64(1)
		withdrawals := []*model.Withdrawal{
			{
				ID:          1,
				Amount:      50.0,
				OrderNumber: "2844830162",
				UserID:      userID,
				CreatedAt:   today,
			},
			{
				ID:          2,
				Amount:      30.0,
				OrderNumber: "28448301623",
				UserID:      userID,
				CreatedAt:   today,
			},
		}

		repo.EXPECT().
			ListWithdrawals(ctx, userID).
			Return(withdrawals, nil)

		resp, err := ctrl.ListWithdrawals(ctx, userID)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp, 2)
	})

	t.Run("Error from repo", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockbalanceRepo(gomock.NewController(t))

		ctrl := New(repo)
		userID := int64(1)

		repo.EXPECT().
			ListWithdrawals(ctx, userID).
			Return(nil, errors.New("error"))

		resp, err := ctrl.ListWithdrawals(ctx, userID)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestController_Withdraw(t *testing.T) {
	t.Parallel()

	t.Run("Withdrawal is created", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockbalanceRepo(gomock.NewController(t))

		ctrl := New(repo)

		userID := int64(1)
		orderNumber := "2844830162"
		amount := 50.0

		req := &model.WithdrawReq{
			Order: orderNumber,
			Sum:   amount,
		}

		repo.EXPECT().
			FindBalanceByUserID(ctx, userID).
			Return(&model.Balance{
				Current:   150.0,
				Withdrawn: 0.0,
			}, nil)

		repo.EXPECT().
			InsertWithdrawal(ctx, userID, orderNumber, amount).
			Return(int64(1), nil)

		err := ctrl.Withdraw(ctx, userID, req)
		require.NoError(t, err)
	})

	t.Run("Not enough balance", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockbalanceRepo(gomock.NewController(t))

		ctrl := New(repo)

		userID := int64(1)
		orderNumber := "2844830162"
		amount := 200.0

		req := &model.WithdrawReq{
			Order: orderNumber,
			Sum:   amount,
		}

		repo.EXPECT().
			FindBalanceByUserID(ctx, userID).
			Return(&model.Balance{
				Current:   150.0,
				Withdrawn: 0.0,
			}, nil)

		err := ctrl.Withdraw(ctx, userID, req)
		require.Error(t, err)
	})

	t.Run("Invalid order number", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockbalanceRepo(gomock.NewController(t))

		ctrl := New(repo)

		userID := int64(1)
		orderNumber := "12345"
		amount := 50.0

		req := &model.WithdrawReq{
			Order: orderNumber,
			Sum:   amount,
		}

		err := ctrl.Withdraw(ctx, userID, req)
		require.Error(t, err)
	})
}
