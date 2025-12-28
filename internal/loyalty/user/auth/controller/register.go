package controller

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/auth/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
	"github.com/yogenyslav/loyalty/pkg/secure"
)

// Register registers a new user and returns its ID.
func (ctrl *Controller) Register(ctx context.Context, req *model.RegisterReq) (*model.RegisterResp, error) {
	hashedPassword, err := secure.HashPassword(req.Password)
	if err != nil {
		return nil, errs.Wrap(err, "hash password")
	}

	user := &model.User{
		Login:          req.Login,
		HashedPassword: hashedPassword,
	}

	userID, err := ctrl.ur.InsertUser(ctx, user)
	if err != nil {
		if errs.CheckUniqueViloation(err) {
			return nil, errs.Wrap(errs.ErrUserAlreadyExists, "insert user")
		}
		return nil, errs.Wrap(err, "insert user")
	}

	accessToken, err := ctrl.jwtProvider.CreateAccessToken(userID)
	if err != nil {
		return nil, errs.Wrap(err, "create access token")
	}

	return &model.RegisterResp{
		ID:    userID,
		Token: accessToken,
	}, nil
}
