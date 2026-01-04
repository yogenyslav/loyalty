package repo

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/auth/model"
	"github.com/yogenyslav/loyalty/pkg/secure"
	db_test "github.com/yogenyslav/loyalty/tests/db"
)

func TestUserRepo(t *testing.T) {
	db := db_test.SetupTestDB(t)
	defer db_test.DropMigrations(t)

	repo := New(db)

	t.Run("user is inserted", func(t *testing.T) {
		defer db_test.ClearTables(t, "loyalty.balance", "loyalty.user")

		ctx := t.Context()
		hashedPassword, err := secure.HashPassword("test123456")
		require.NoError(t, err)

		u := &model.User{
			Login:          "test_user",
			HashedPassword: hashedPassword,
		}
		userID, err := repo.InsertUser(ctx, u)
		require.NoError(t, err)
		require.Greater(t, userID, int64(0))

		foundUser, err := repo.FindUserByLogin(ctx, "test_user")
		require.NoError(t, err)
		require.Equal(t, u.Login, foundUser.Login)
		require.Equal(t, u.HashedPassword, foundUser.HashedPassword)
	})

	t.Run("Insert duplicate user", func(t *testing.T) {
		defer db_test.ClearTables(t, "loyalty.balance", "loyalty.user")

		ctx := t.Context()
		hashedPassword, err := secure.HashPassword("test123456")
		require.NoError(t, err)

		u := &model.User{
			Login:          "test_user",
			HashedPassword: hashedPassword,
		}
		userID, err := repo.InsertUser(ctx, u)
		require.NoError(t, err)
		require.Greater(t, userID, int64(0))

		_, err = repo.InsertUser(ctx, u)
		require.Error(t, err)
	})
}
