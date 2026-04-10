package domain

import "time"

type User struct {
    ID           uint   `gorm:"primaryKey"`
    Username     string `gorm:"uniqueIndex;not null"`
    PasswordHash string `gorm:"column:password_hash;not null"`
    IsAdmin      bool   `gorm:"column:is_admin;default:false"`
    CreatedAt    time.Time
}