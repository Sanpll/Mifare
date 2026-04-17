package domain

import "github.com/shopspring/decimal"

type Card struct {
	ID         uint            `db:"id"          json:"id"`
	CardNumber string          `db:"card_number" json:"card_number"`
	Balance    decimal.Decimal `db:"balance"     json:"balance"`
	IsBlocked  bool            `db:"is_blocked"  json:"is_blocked"`
	OwnerName  string          `db:"owner_name"  json:"ownerName"`
	KeyID      uint            `db:"key_id"      json:"key_id"`
}

type CardWithKey struct {
	ID         uint            `db:"id"`
	CardNumber string          `db:"card_number"`
	Balance    decimal.Decimal `db:"balance"`
	IsBlocked  bool            `db:"is_blocked"`
	OwnerName  string          `db:"owner_name"`
	KeyID      uint            `db:"key_id"`
	KeyValue   string          `db:"key_value"`
}