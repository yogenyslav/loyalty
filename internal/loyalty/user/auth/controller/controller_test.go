package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/auth/model"
	"github.com/yogenyslav/loyalty/pkg/secure"
	"github.com/yogenyslav/loyalty/tests/mocks"
	"go.uber.org/mock/gomock"
)

func TestController_Login(t *testing.T) {
	t.Parallel()

	t.Run("User is found, login success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockuserRepo(gomock.NewController(t))
		jwt := mocks.NewMockJwtProvider(gomock.NewController(t))

		ctrl := New(repo, jwt)

		req := &model.LoginReq{
			Login:    "test_user",
			Password: "test123456",
		}

		hashedPassword, err := secure.HashPassword(req.Password)
		require.NoError(t, err)

		user := &model.User{
			ID:             1,
			Login:          "test_user",
			HashedPassword: hashedPassword,
		}

		repo.EXPECT().
			FindUserByLogin(ctx, req.Login).
			Return(user, nil)

		jwt.EXPECT().
			CreateAccessToken(user.ID).
			Return("token", nil)

		resp, err := ctrl.Login(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, model.LoginResp{Token: "token"}, *resp)
	})

	t.Run("User not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockuserRepo(gomock.NewController(t))
		jwt := mocks.NewMockJwtProvider(gomock.NewController(t))

		ctrl := New(repo, jwt)

		req := &model.LoginReq{
			Login:    "asdf",
			Password: "test123456",
		}

		repo.EXPECT().
			FindUserByLogin(ctx, req.Login).
			Return(nil, pgx.ErrNoRows)

		resp, err := ctrl.Login(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("Invalid password", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockuserRepo(gomock.NewController(t))
		jwt := mocks.NewMockJwtProvider(gomock.NewController(t))

		ctrl := New(repo, jwt)

		req := &model.LoginReq{
			Login:    "test_user",
			Password: "asdf",
		}

		hashedPassword, err := secure.HashPassword("test123456")
		require.NoError(t, err)

		user := &model.User{
			ID:             1,
			Login:          "test_user",
			HashedPassword: hashedPassword,
		}

		repo.EXPECT().
			FindUserByLogin(ctx, req.Login).
			Return(user, nil)

		resp, err := ctrl.Login(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("SQL error on FindUserByLogin", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockuserRepo(gomock.NewController(t))
		jwt := mocks.NewMockJwtProvider(gomock.NewController(t))

		ctrl := New(repo, jwt)

		req := &model.LoginReq{
			Login:    "test_user",
			Password: "test123456",
		}

		repo.EXPECT().
			FindUserByLogin(ctx, req.Login).
			Return(nil, errors.New("some error"))

		resp, err := ctrl.Login(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("Error on CreateAccessToken", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockuserRepo(gomock.NewController(t))
		jwt := mocks.NewMockJwtProvider(gomock.NewController(t))

		ctrl := New(repo, jwt)

		req := &model.LoginReq{
			Login:    "test_user",
			Password: "test123456",
		}

		hashedPassword, err := secure.HashPassword(req.Password)
		require.NoError(t, err)

		user := &model.User{
			ID:             1,
			Login:          "test_user",
			HashedPassword: hashedPassword,
		}

		repo.EXPECT().
			FindUserByLogin(ctx, req.Login).
			Return(user, nil)

		jwt.EXPECT().
			CreateAccessToken(user.ID).
			Return("", errors.New("some error"))

		resp, err := ctrl.Login(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestController_Register(t *testing.T) {
	t.Parallel()

	t.Run("User is created", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockuserRepo(gomock.NewController(t))
		jwt := mocks.NewMockJwtProvider(gomock.NewController(t))

		ctrl := New(repo, jwt)

		req := &model.RegisterReq{
			Login:    "new_user",
			Password: "test123456",
		}

		repo.EXPECT().
			InsertUser(gomock.Any(), gomock.Any()).
			Return(int64(1), nil)

		jwt.EXPECT().
			CreateAccessToken(int64(1)).
			Return("token", nil)

		resp, err := ctrl.Register(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, model.RegisterResp{
			ID:    1,
			Token: "token",
		}, *resp)
	})

	t.Run("User already exists", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockuserRepo(gomock.NewController(t))
		jwt := mocks.NewMockJwtProvider(gomock.NewController(t))

		ctrl := New(repo, jwt)

		req := &model.RegisterReq{
			Login:    "test_user",
			Password: "test123456",
		}

		repo.EXPECT().
			InsertUser(gomock.Any(), gomock.Any()).
			Return(int64(0), &pgconn.PgError{Code: pgerrcode.UniqueViolation})

		resp, err := ctrl.Register(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("SQL error on InsertUser", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockuserRepo(gomock.NewController(t))
		jwt := mocks.NewMockJwtProvider(gomock.NewController(t))

		ctrl := New(repo, jwt)

		req := &model.RegisterReq{
			Login:    "test_user",
			Password: "test123456",
		}

		repo.EXPECT().
			InsertUser(gomock.Any(), gomock.Any()).
			Return(int64(0), errors.New("some error"))

		resp, err := ctrl.Register(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("Error on CreateAccessToken", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockuserRepo(gomock.NewController(t))
		jwt := mocks.NewMockJwtProvider(gomock.NewController(t))

		ctrl := New(repo, jwt)

		req := &model.RegisterReq{
			Login:    "test_user",
			Password: "test123456",
		}

		repo.EXPECT().
			InsertUser(gomock.Any(), gomock.Any()).
			Return(int64(1), nil)

		jwt.EXPECT().
			CreateAccessToken(int64(1)).
			Return("", errors.New("some error"))

		resp, err := ctrl.Register(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("Error hashing password", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := mocks.NewMockuserRepo(gomock.NewController(t))
		jwt := mocks.NewMockJwtProvider(gomock.NewController(t))

		ctrl := New(repo, jwt)

		req := &model.RegisterReq{
			Login:    "test_user",
			Password: string(make([]byte, 1000)),
		}

		resp, err := ctrl.Register(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}
