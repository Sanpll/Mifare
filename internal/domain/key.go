package domain

type Key struct {
	ID          uint   `gorm:"primaryKey"`
	KeyValue    string `gorm:"column:key_value;not null"`
	KeyType     string `gorm:"column:key_type"`
	Description string
	Cards       []Card `gorm:"foreignKey:KeyID"`
}
