package repository

import (
	"fmt"
	"mifare/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user domain.User) (int, error) {
	query := `
        INSERT INTO users (username, password_hash, is_admin)
        VALUES (?, ?, ?)
        RETURNING id`

	var id int
	err := r.db.QueryRowx(query, user.Username, user.PasswordHash, user.IsAdmin).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (r *UserRepository) GetByUsername(username string) (domain.User, error) {
	query := `
        SELECT id, username, password_hash, is_admin, created_at
        FROM users
        WHERE username = ?`

	var user domain.User
	err := r.db.Get(&user, query, username)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetById(id int) (domain.User, error) {
	query := `
        SELECT id, username, is_admin, created_at
        FROM users
        WHERE id = ?`

	var user domain.User
	err := r.db.Get(&user, query, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetAll() ([]domain.User, error) {
	query := `
        SELECT id, username, is_admin, created_at
        FROM users`

	var users []domain.User
	err := r.db.Select(&users, query)
	if err != nil {
		return []domain.User{}, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (r *UserRepository) Update(id int, user domain.User) error {
	query := `
        UPDATE users 
        SET username = ? 
        WHERE id = ?`

	result, err := r.db.Exec(query, user.Username, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

func (r *UserRepository) Delete(id int) error {
	query := `
        DELETE FROM users
        WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}
