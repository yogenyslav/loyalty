// Package model contains data models related to balances.
package model

// Balance represents a user's loyalty balance.
type Balance struct {
	UserID    string  `db:"fk_user_id"`
	Current   float64 `db:"current"`
	Withdrawn float64 `db:"withdrawn"`
	UpdatedAt int64   `db:"updated_at"`
}

// ToDto converts Balance to BalanceDto.
func (b *Balance) ToDto() *BalanceDto {
	return &BalanceDto{
		Current:   b.Current,
		Withdrawn: b.Withdrawn,
	}
}

// BalanceDto is a data transfer object for Balance.
type BalanceDto struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
