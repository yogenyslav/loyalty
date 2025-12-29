package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

func (h *Handler) ListWithdrawals(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return errs.Wrap(errs.ErrUnauthorized, "get userID from context")
	}

	withdrawals, err := h.ctrl.ListWithdrawals(c.Context(), userID)
	if err != nil {
		return errs.Wrap(err, "list withdrawals")
	}

	if len(withdrawals) == 0 {
		return c.SendStatus(http.StatusNoContent)
	}

	return c.JSON(withdrawals)
}
