package repository

import (
	"fmt"
	"mifare/internal/domain"

	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthSQLite(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) CreateUser(user domain.User) (int, error) {
	query := `
        INSERT INTO users (username, password_hash, is_admin)
        VALUES (:username, :password_hash, :is_admin)
        RETURNING id`

    var id int
    if err := r.db.Get(&id, query, user); err != nil {
        return 0, fmt.Errorf("failed to create user: %w", err)
    }

    return id, nil
}
