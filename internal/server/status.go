package server

import (
	"net/http"

	jwtwire "github.com/golang-jwt/jwt/v5"
	"github.com/yogenyslav/loyalty/pkg/errs"
	"github.com/yogenyslav/loyalty/pkg/jwt"
	"github.com/yogenyslav/loyalty/pkg/luhn"
)

var errStatus = map[error]ErrorResponse{ //nolint:gochecknoglobals // used for mapping errors to status codes.
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
	errs.ErrUnauthorized: {
		Status: http.StatusUnauthorized,
	},
	jwtwire.ErrTokenExpired: {
		ErrMessage: "access token has expired",
		Status:     http.StatusUnauthorized,
	},
	// 402 Payment Required
	errs.ErrNotEnoughBalance: {
		Status: http.StatusPaymentRequired,
	},
	// 409 Conflict
	errs.ErrUserAlreadyExists: {
		Status: http.StatusConflict,
	},
	errs.ErrOrderCreatedByAnotherUser: {
		Status: http.StatusConflict,
	},
	// 422 Unprocessable Entity
	luhn.ErrInvalidNumber: {
		ErrMessage: "order number failed Luhn check",
		Status:     http.StatusUnprocessableEntity,
	},
}
