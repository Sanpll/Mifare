package domain

type Card struct {
	ID           uint   `gorm:"primaryKey"`
	CardNumber   string `gorm:"uniqueIndex;not null"`
	Balance      float64
	IsBlocked    bool
	OwnerName    string
	KeyID        uint
	Key          Key           `gorm:"foreignKey:KeyID"`
	Transactions []Transaction `gorm:"foreignKey:CardID"`
}
