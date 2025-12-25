package controller

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/auth/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
	"github.com/yogenyslav/loyalty/pkg/secure"
)

// Login authenticates a user and returns a JWT token upon successful authentication.
func (ctrl *Controller) Login(ctx context.Context, req model.LoginReq) (model.LoginResp, error) {
	var resp model.LoginResp

	userDB, err := ctrl.ur.FindUserByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return resp, errs.Wrap(errs.ErrInvalidCredentials, "user not found")
		}
		return resp, errs.Wrap(err, "find user by login")
	}

	if !secure.VerifyPassword(userDB.HashedPassword, req.Password) {
		return resp, errs.Wrap(errs.ErrInvalidCredentials, "verify password")
	}

	resp.Token, err = ctrl.jwtProvider.CreateAccessToken(userDB.ID)
	return resp, errs.Wrap(err, "create access token")
}
