package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/auth/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

// Login handles login HTTP requests.
func (h *Handler) Login(c fiber.Ctx) error {
	var req model.LoginReq
	if err := c.Bind().JSON(&req); err != nil {
		return errs.Wrap(errs.ErrInvalidRequest, "unmarshal req")
	}

	resp, err := h.ac.Login(c.Context(), &req)
	if err != nil {
		return errs.Wrap(err, "user login")
	}

	return c.JSON(resp)
}
