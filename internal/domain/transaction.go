package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID         uint            `db:"id"          json:"id"`
	Price      decimal.Decimal `db:"price"       json:"price"`
	CardID     uint            `db:"card_id"     json:"card_id"`
	TerminalID uint            `db:"terminal_id" json:"terminal_id"`
	Status     string          `db:"status"      json:"status"`
	CreatedAt  time.Time       `db:"created_at"  json:"created_at"`
}
