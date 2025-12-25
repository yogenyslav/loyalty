package server

import (
	"errors"
	"maps"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// ErrorResponse is a struct that holds the error message and status code.
type ErrorResponse struct {
	Msg    string `json:"msg"`
	Status int    `json:"-"`
}

// ErrorHandler is a struct that holds the error status map.
type ErrorHandler struct {
	status map[error]ErrorResponse
}

// NewErrorHandler creates a new ErrorHandler instance with the given error status map.
func NewErrorHandler(errStatus map[error]ErrorResponse) ErrorHandler {
	status := map[error]ErrorResponse{
		pgx.ErrNoRows: {
			Msg:    "no requested resource",
			Status: http.StatusNotFound,
		},
		fiber.ErrNotFound: {
			Msg:    "route not found",
			Status: http.StatusNotFound,
		},
		fiber.ErrUnprocessableEntity: {
			Msg:    "validation error",
			Status: http.StatusUnprocessableEntity,
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
	log.Err(err).Msg(e.Msg)
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
			break
		}
	}

	if !ok {
		e = ErrorResponse{
			Msg:    "unknown error",
			Status: http.StatusInternalServerError,
		}
	}
	if e.Msg == "" {
		e.Msg = err.Error()
	}

	return e
}
