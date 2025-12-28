package accrual

import "github.com/yogenyslav/loyalty/internal/loyalty/user/orders/model"

// AccrualInfo represents the accrual information for an order.
type AccrualInfo struct {
	Order   string       `json:"order"`
	Accrual float64      `json:"accrual"`
	Status  model.Status `json:"status"`
}
