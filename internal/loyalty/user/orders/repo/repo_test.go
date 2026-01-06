package repo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/loyalty/internal/accrual"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"
	db_test "github.com/yogenyslav/loyalty/tests/db"
)

func TestOrdersRepo(t *testing.T) {
	time.Sleep(time.Second)
	db := db_test.SetupTestDB(t)
	defer db_test.DropMigrations(t)

	repo := New(db)

	t.Run("Order is created", func(t *testing.T) {
		defer db_test.ClearTables(t, "loyalty.order", "loyalty.user")

		ctx := t.Context()
		userID := db_test.InsertUser(t, db)

		order := &model.Order{
			Number: "2844830162",
			UserID: userID,
		}
		err := repo.InsertOrder(ctx, order)
		require.NoError(t, err)

		orderFromDB, err := repo.FindOrderByNumber(ctx, order.Number)
		require.NoError(t, err)
		require.Equal(t, order.Number, orderFromDB.Number)
		require.Equal(t, order.UserID, orderFromDB.UserID)
	})

	t.Run("Get orders list", func(t *testing.T) {
		defer db_test.ClearTables(t, "loyalty.order", "loyalty.user")

		ctx := t.Context()
		userID := db_test.InsertUser(t, db)

		ordersToInsert := []*model.Order{
			{Number: "2844830162", UserID: userID},
			{Number: "28448301623", UserID: userID},
		}

		for _, order := range ordersToInsert {
			err := repo.InsertOrder(ctx, order)
			require.NoError(t, err)
		}

		ordersFromDB, err := repo.ListOrders(ctx, userID)
		require.NoError(t, err)
		require.Len(t, ordersFromDB, len(ordersToInsert))
	})

	t.Run("Update order", func(t *testing.T) {
		defer db_test.ClearTables(t, "loyalty.order", "loyalty.user")

		ctx := t.Context()
		userID := db_test.InsertUser(t, db)

		order := &model.Order{
			Number: "2844830162",
			UserID: userID,
		}
		err := repo.InsertOrder(ctx, order)
		require.NoError(t, err)

		data := &accrual.OrderAccrual{
			Order:   order.Number,
			Status:  "PROCESSED",
			Accrual: 150.0,
		}
		err = repo.UpdateOrderAfterPolling(ctx, data)
		require.NoError(t, err)

		updatedOrder, err := repo.FindOrderByNumber(ctx, order.Number)
		require.NoError(t, err)
		require.Equal(t, data.Status, updatedOrder.Status)
		require.Equal(t, data.Accrual, updatedOrder.Accrual)
	})
}

func TestAccrualPolling(t *testing.T) {
	time.Sleep(time.Second)
	db := db_test.SetupTestDB(t)
	defer db_test.DropMigrations(t)

	repo := New(db)

	t.Run("Schedule and find pollable orders", func(t *testing.T) {
		defer db_test.ClearTables(t, "loyalty.poll_accrual", "loyalty.order", "loyalty.user")

		ctx := t.Context()
		userID := db_test.InsertUser(t, db)

		order := &model.Order{
			Number: "2844830162",
			UserID: userID,
		}
		err := repo.InsertOrder(ctx, order)
		require.NoError(t, err)

		err = repo.SchedulePolling(ctx, order.Number)
		require.NoError(t, err)

		pollableOrders, err := repo.FindPollableOrders(ctx)
		require.NoError(t, err)
		require.Len(t, pollableOrders, 1)
		_, ok := pollableOrders[order.Number]
		require.True(t, ok)
	})

	t.Run("Update polling schedule", func(t *testing.T) {
		defer db_test.ClearTables(t, "loyalty.poll_accrual", "loyalty.order", "loyalty.user")

		ctx := t.Context()
		userID := db_test.InsertUser(t, db)

		order := &model.Order{
			Number: "2844830162",
			UserID: userID,
		}
		err := repo.InsertOrder(ctx, order)
		require.NoError(t, err)

		err = repo.SchedulePolling(ctx, order.Number)
		require.NoError(t, err)

		err = repo.UpdatePollingSchedule(ctx, order.Number, 1)
		require.NoError(t, err)

		time.Sleep(time.Second)

		pollableOrders, err := repo.FindPollableOrders(ctx)
		require.NoError(t, err)
		require.Len(t, pollableOrders, 1)
	})

	t.Run("Delete polling schedule", func(t *testing.T) {
		defer db_test.ClearTables(t, "loyalty.poll_accrual", "loyalty.order", "loyalty.user")

		ctx := t.Context()
		userID := db_test.InsertUser(t, db)

		order := &model.Order{
			Number: "2844830162",
			UserID: userID,
		}
		err := repo.InsertOrder(ctx, order)
		require.NoError(t, err)

		err = repo.SchedulePolling(ctx, order.Number)
		require.NoError(t, err)

		err = repo.DeletePollingSchedule(ctx, order.Number)
		require.NoError(t, err)

		pollableOrders, err := repo.FindPollableOrders(ctx)
		require.NoError(t, err)
		require.Len(t, pollableOrders, 0)
	})
}
