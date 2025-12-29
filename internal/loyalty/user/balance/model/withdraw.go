package model

import "time"

// Withdrawal represents a withdrawal record of loyalty points.
type Withdrawal struct {
	ID          int64     `db:"id"`
	Amount      float64   `db:"amount"`
	OrderNumber string    `db:"order_number"`
	UserID      int64     `db:"fk_user_id"`
	CreatedAt   time.Time `db:"created_at"`
}

// ToDto converts Withdrawal to WithdrawalDto.
func (w *Withdrawal) ToDto() *WithdrawalDto {
	return &WithdrawalDto{
		Order:       w.OrderNumber,
		Sum:         w.Amount,
		ProcessedAt: w.CreatedAt.Format(time.RFC3339),
	}
}

// WithdrawalDto is a data transfer object for Withdrawal.
type WithdrawalDto struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

// WithdrawReq is a request model for withdrawing loyalty points.
type WithdrawReq struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
