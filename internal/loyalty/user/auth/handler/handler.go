// Package handler routes HTTP requests for user authentication.
package handler

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/auth/model"
)

type authController interface {
	Register(ctx context.Context, req model.RegisterReq) (model.RegisterResp, error)
	Login(ctx context.Context, req model.LoginReq) (model.LoginResp, error)
}

// Handler provides methods to handle HTTP requests for user authentication.
type Handler struct {
	ac authController
}

// New creates a new Handler instance.
func New(ac authController) *Handler {
	return &Handler{ac: ac}
}
