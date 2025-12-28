package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

// ListOrders handles the HTTP request to list user orders.
func (h *Handler) ListOrders(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return errs.Wrap(errs.ErrUnauthorized, "get userID from context")
	}

	orders, err := h.oc.ListOrders(c.Context(), userID)
	if err != nil {
		return errs.Wrap(err, "list orders")
	}

	if len(orders) == 0 {
		return c.SendStatus(http.StatusNoContent)
	}
	return c.JSON(orders)
}
