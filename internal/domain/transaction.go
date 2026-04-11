package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID         uint            `db:"id" json:"id"`
	Price      decimal.Decimal `db:"price" json:"price"`
	CardID     uint            `db:"card_id" json:"cardId"`
	TerminalID uint            `db:"terminal_id" json:"terminalId"`
	Status     string          `db:"status" json:"status"`
	CreatedAt  time.Time       `db:"created_at" json:"createdAt"`
}
