package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID         uint `gorm:"primaryKey"`
	Price      decimal.Decimal
	CardID     uint
	Card       Card `gorm:"foreignKey:CardID"`
	TerminalID uint
	Terminal   Terminal `gorm:"foreignKey:TerminalID"`
	Status     string
	CreatedAt  time.Time
}
