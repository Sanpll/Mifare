package domain

type Terminal struct {
	ID           uint   `gorm:"primaryKey"`
	SerialNumber string `gorm:"uniqueIndex;not null"`
	Address      string
	Name         string
	Transactions []Transaction `gorm:"foreignKey:TerminalID"`
}
