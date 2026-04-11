package domain

import "time"

type User struct {
	ID           uint      `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	IsAdmin      bool      `db:"is_admin" json:"isAdmin"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
}
