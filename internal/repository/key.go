package repository

import (
	"fmt"
	"mifare/internal/domain"
	"strings"

	"github.com/jmoiron/sqlx"
)

type KeyRepository struct {
	db *sqlx.DB
}

func NewKeyRepository(db *sqlx.DB) *KeyRepository {
	return &KeyRepository{
		db: db,
	}
}

func (r *KeyRepository) Create(key domain.Key) (int, error) {
	query := `
        INSERT INTO keys (key_value, key_type, description)
        VALUES (?, ?, ?)
        RETURNING id`

	var id int
	err := r.db.QueryRowx(query, key.KeyValue, key.KeyType, key.Description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create key: %w", err)
	}

	return id, nil
}

func (r *KeyRepository) GetById(id int) (domain.Key, error) {
	query := `
        SELECT id, key_value, key_type, description
        FROM keys
        WHERE id = ?`

	var key domain.Key
	err := r.db.Get(&key, query, id)
	if err != nil {
		return domain.Key{}, fmt.Errorf("failed to get key: %w", err)
	}

	return key, nil
}

func (r *KeyRepository) GetAll() ([]domain.Key, error) {
	query := `
        SELECT id, key_value, key_type, description
        FROM keys`

	var keys []domain.Key
	err := r.db.Select(&keys, query)
	if err != nil {
		return []domain.Key{}, fmt.Errorf("failed to get keys: %w", err)
	}

	return keys, nil
}

func (r *KeyRepository) GetByValue(keyValue string) (domain.Key, error) {
	if keyValue == "" {
		return domain.Key{}, fmt.Errorf("key_value is required")
	}

	query := `
		SELECT id, key_value, key_type, description
		FROM keys
		WHERE key_value = ?`

	var key domain.Key
	err := r.db.Get(&key, query, keyValue)
	if err != nil {
		return domain.Key{}, fmt.Errorf("failed to get key id by value: %w", err)
	}

	return key, nil
}

func (r *KeyRepository) Update(id int, key domain.Key) error {
	var updates []string
	var args []interface{}

	if key.KeyValue != "" {
		updates = append(updates, "key_value = ?")
		args = append(args, key.KeyValue)
	}
	if key.KeyType != "" {
		updates = append(updates, "key_type = ?")
		args = append(args, key.KeyType)
	}
	if key.Description != "" {
		updates = append(updates, "description = ?")
		args = append(args, key.Description)
	}

	if len(updates) == 0 {
		return fmt.Errorf("you must update at least 1 param")
	}

	query := fmt.Sprintf(`
        UPDATE keys 
        SET %s 
        WHERE id = ?`, 
		strings.Join(updates, ", "))

	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("key with id %d not found", id)
	}

	return nil
}

func (r *KeyRepository) Delete(id int) error {
	query := `
        DELETE FROM keys
        WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("key with id %d not found", id)
	}

	return nil
}
