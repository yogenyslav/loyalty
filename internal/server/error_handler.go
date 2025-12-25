package server

import (
	"errors"
	"maps"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

// ErrorResponse is a struct that holds the error message and status code.
type ErrorResponse struct {
	ErrMessage string `json:"err_message"`
	Status     int    `json:"-"`
}

// ErrorHandler is a struct that holds the error status map.
type ErrorHandler struct {
	status map[error]ErrorResponse
}

// NewErrorHandler creates a new ErrorHandler instance with the given error status map.
func NewErrorHandler(errStatus map[error]ErrorResponse) ErrorHandler {
	status := map[error]ErrorResponse{
		pgx.ErrNoRows: {
			ErrMessage: "no requested resource",
			Status:     http.StatusNotFound,
		},
		fiber.ErrNotFound: {
			ErrMessage: "route not found",
			Status:     http.StatusNotFound,
		},
		fiber.ErrUnprocessableEntity: {
			ErrMessage: "validation error",
			Status:     http.StatusUnprocessableEntity,
		},
	}

	maps.Copy(status, errStatus)

	return ErrorHandler{
		status: status,
	}
}

// Handler is a method that handles the error and returns a JSON response.
func (h ErrorHandler) Handler(c fiber.Ctx, err error) error {
	e := h.getErrorResponse(err)
	return c.Status(e.Status).JSON(e)
}

func (h ErrorHandler) getErrorResponse(err error) ErrorResponse {
	var (
		ok bool
		e  ErrorResponse
	)

	for k, v := range h.status {
		if errors.Is(err, k) {
			ok = true
			e = v

			if e.ErrMessage == "" {
				e.ErrMessage = k.Error()
			}

			break
		}
	}

	if !ok {
		e = ErrorResponse{
			ErrMessage: "unknown error",
			Status:     http.StatusInternalServerError,
		}
	}

	return e
}
