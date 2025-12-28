package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

// ProcessOrder handles the HTTP request to process a new order.
func (h *Handler) ProcessOrder(c fiber.Ctx) error {
	orderNumber := string(c.Body())

	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return errs.Wrap(errs.ErrUnauthorized, "get userID from context")
	}

	processedOrderID, err := h.oc.ProcessOrder(c.Context(), orderNumber, userID)
	if err != nil {
		return errs.Wrap(err, "process order")
	}

	if processedOrderID == "" {
		return c.SendStatus(http.StatusOK)
	}
	return c.SendStatus(http.StatusAccepted)
}
