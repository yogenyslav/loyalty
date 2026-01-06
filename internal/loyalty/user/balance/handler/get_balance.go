package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

// GetBalance handles the requests to get user's balance.
func (h *Handler) GetBalance(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return errs.Wrap(errs.ErrUnauthorized, "get userID from context")
	}

	balance, err := h.ctrl.GetBalance(c.Context(), userID)
	if err != nil {
		return errs.Wrap(err, "get balance")
	}

	return c.Status(fiber.StatusOK).JSON(balance)
}
