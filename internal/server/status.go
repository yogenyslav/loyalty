package server

import (
	"net/http"

	"github.com/yogenyslav/loyalty/pkg/errs"
)

var errStatus = map[error]ErrorResponse{
	errs.ErrInvalidRequest: {
		Status: http.StatusBadRequest,
	},
	errs.ErrInvalidCredentials: {
		Status: http.StatusUnauthorized,
	},
	errs.ErrUserAlreadyExists: {
		Status: http.StatusConflict,
	},
}
