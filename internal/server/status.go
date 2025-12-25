package server

import (
	"net/http"

	"github.com/yogenyslav/loyalty/pkg/errs"
	"github.com/yogenyslav/loyalty/pkg/jwt"
)

var errStatus = map[error]ErrorResponse{
	// 400 Bad Request
	errs.ErrInvalidRequest: {
		Status: http.StatusBadRequest,
	},
	// 401 Unauthorized
	errs.ErrInvalidCredentials: {
		Status: http.StatusUnauthorized,
	},
	jwt.ErrJwtSignMethod: {
		Status: http.StatusUnauthorized,
	},
	errs.ErrInvalidToken: {
		Status: http.StatusUnauthorized,
	},
	errs.ErrMisingAccessToken: {
		Status: http.StatusUnauthorized,
	},
	// 409 Conflict
	errs.ErrUserAlreadyExists: {
		Status: http.StatusConflict,
	},
}
