package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/auth/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

// Register handles registration HTTP requests.
func (h *Handler) Register(c fiber.Ctx) error {
	var req model.RegisterReq
	if err := c.Bind().JSON(&req); err != nil {
		return errs.Wrap(errs.ErrInvalidRequest, "unmarshal req")
	}

	resp, err := h.ac.Register(c.Context(), &req)
	if err != nil {
		return errs.Wrap(err, "user register")
	}

	return c.JSON(resp)
}
