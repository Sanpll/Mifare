package dto

import "time"

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type UsersResponse struct {
	Users []UserResponse `json:"users"`
}

type UserUpdate struct {
	Username string `json:"username" binding:"omitempty"`
}