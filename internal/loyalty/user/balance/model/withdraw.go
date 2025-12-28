package model

// WithdrawReq is a request model for withdrawing loyalty points.
type WithdrawReq struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
