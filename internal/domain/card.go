package domain

import "github.com/shopspring/decimal"

type Card struct {
	ID         uint            `db:"id" json:"id"`
	CardNumber string          `db:"card_number" json:"cardNumber"`
	Balance    decimal.Decimal `db:"balance" json:"balance"`
	IsBlocked  bool            `db:"is_blocked" json:"isBlocked"`
	OwnerName  string          `db:"owner_name" json:"ownerName"`
	KeyID      uint            `db:"key_id" json:"keyId"`
}
