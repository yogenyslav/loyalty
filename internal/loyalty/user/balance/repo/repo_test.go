package repo

import (
	"testing"

	"github.com/stretchr/testify/require"
	db_test "github.com/yogenyslav/loyalty/tests/db"
)

func TestBalanceRepo(t *testing.T) {
	db := db_test.SetupTestDB(t)
	defer db_test.DropMigrations(t)

	repo := New(db)

	t.Run("Balance is created", func(t *testing.T) {
		defer db_test.ClearTables(t, "loyalty.balance", "loyalty.user")

		ctx := t.Context()
		userID := db_test.InsertUser(t, db)

		err := repo.InsertBalance(ctx, userID)
		require.NoError(t, err)
	})

	t.Run("Withdrawal is created and balance updated", func(t *testing.T) {
		defer db_test.ClearTables(t, "loyalty.withdrawal", "loyalty.balance", "loyalty.user")

		ctx := t.Context()
		userID := db_test.InsertUser(t, db)

		err := repo.InsertBalance(ctx, userID)
		require.NoError(t, err)

		err = repo.UpdateBalanceAccrual(ctx, userID, 100.0)
		require.NoError(t, err)

		withdrawalID, err := repo.InsertWithdrawal(ctx, userID, "2844830162", 50.0)
		require.NoError(t, err)
		require.Greater(t, withdrawalID, int64(0))

		err = repo.UpdateBalanceWithdraw(ctx, userID, 50.0)
		require.NoError(t, err)

		updatedBalance, err := repo.FindBalanceByUserID(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, 50.0, updatedBalance.Current)

		withdrawals, err := repo.ListWithdrawals(ctx, userID)
		require.NoError(t, err)
		require.Len(t, withdrawals, 1)
		require.Equal(t, "2844830162", withdrawals[0].OrderNumber)
		require.Equal(t, 50.0, withdrawals[0].Amount)
		require.Equal(t, userID, withdrawals[0].UserID)
	})
}
