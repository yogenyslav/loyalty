package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/yogenyslav/loyalty/internal/loyalty/user/balance/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

// Withdraw handles the requests to withdraw user's balance.
func (h *Handler) Withdraw(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return errs.Wrap(errs.ErrUnauthorized, "get userID from context")
	}

	var req model.WithdrawReq
	if err := c.Bind().JSON(&req); err != nil {
		return errs.Wrap(errs.ErrInvalidRequest, "unmarhsal request")
	}

	if err := h.ctrl.Withdraw(c.Context(), userID, &req); err != nil {
		return errs.Wrap(err, "withdraw")
	}

	return c.SendStatus(fiber.StatusOK)
}
