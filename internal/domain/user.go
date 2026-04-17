package domain

import "time"

type User struct {
	ID           uint      `db:"id"            json:"id"`
	Username     string    `db:"username"      json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	IsAdmin      bool      `db:"is_admin"      json:"is_admin"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
}
