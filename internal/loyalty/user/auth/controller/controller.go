// Package controller implements logic layer for user authentication.
package controller

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/auth/model"
	"github.com/yogenyslav/loyalty/pkg/database"
)

//go:generate mockgen -destination=../../../../../tests/mocks/user_repo.go -package=mocks . userRepo,userBalanceRepo
type userRepo interface {
	InsertUser(ctx context.Context, u *model.User) (int64, error)
	FindUserByLogin(ctx context.Context, login string) (*model.User, error)
}

type userBalanceRepo interface {
	InsertBalance(ctx context.Context, userID int64) error
}

// JwtProvider defines the interface for JWT token generation.
//
//go:generate mockgen -destination=../../../../../tests/mocks/jwt.go -package=mocks . JwtProvider
type JwtProvider interface {
	CreateAccessToken(userID int64) (string, error)
}

// Controller provides methods for user authentication logic.
type Controller struct {
	ur          userRepo
	br          userBalanceRepo
	jwtProvider JwtProvider
	uow         database.UnitOfWork
}

// New creates a new Controller instance.
func New(ur userRepo, br userBalanceRepo, jwtProvider JwtProvider, uow database.UnitOfWork) *Controller {
	return &Controller{
		ur:          ur,
		br:          br,
		jwtProvider: jwtProvider,
		uow:         uow,
	}
}
