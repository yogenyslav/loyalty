// Package controller implements logic layer for user authentication.
package controller

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/auth/model"
)

//go:generate mockgen -destination=../../../../../tests/mocks/user_repo.go -package=mocks . userRepo
type userRepo interface {
	InsertUser(ctx context.Context, u *model.User) (int64, error)
	FindUserByLogin(ctx context.Context, login string) (*model.User, error)
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
	jwtProvider JwtProvider
}

// New creates a new Controller instance.
func New(ur userRepo, jwtProvider JwtProvider) *Controller {
	return &Controller{
		ur:          ur,
		jwtProvider: jwtProvider,
	}
}
