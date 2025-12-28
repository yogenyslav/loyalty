// Package model contains data models related to orders.
package model

import "time"

// Order represents an order.
type Order struct {
	Number    string    `db:"number"`
	Status    Status    `db:"status"`
	Accrual   float64   `db:"accrual"`
	UserID    int64     `db:"fk_user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// ToDto converts Order to OrderDto.
func (o *Order) ToDto() *OrderDto {
	return &OrderDto{
		Number:     o.Number,
		Status:     o.Status,
		Accrual:    o.Accrual,
		UploadedAt: o.CreatedAt.Format(time.RFC3339),
	}
}

// OrderDto is a data transfer object for Order.
type OrderDto struct {
	Number     string  `json:"number"`
	Status     Status  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}
