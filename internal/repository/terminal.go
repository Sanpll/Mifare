package repository

import (
	"fmt"
	"mifare/internal/domain"
	"strings"

	"github.com/jmoiron/sqlx"
)

type TerminalRepository struct {
	db *sqlx.DB
}

func NewTerminalRepository(db *sqlx.DB) *TerminalRepository {
	return &TerminalRepository{
		db: db,
	}
}

func (r *TerminalRepository) Create(terminal domain.Terminal) (int, error) {
	query := `
        INSERT INTO terminals (serial_number, address, name)
        VALUES (?, ?, ?)
        RETURNING id`

	var id int
	err := r.db.QueryRowx(query, terminal.SerialNumber, terminal.Address, terminal.Name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create terminal: %w", err)
	}

	return id, nil
}

func (r *TerminalRepository) GetById(id int) (domain.Terminal, error) {
	query := `
        SELECT id, serial_number, address, name
        FROM terminals
        WHERE id = ?`

	var terminal domain.Terminal
	err := r.db.Get(&terminal, query, id)
	if err != nil {
		return domain.Terminal{}, fmt.Errorf("failed to get terminal: %w", err)
	}

	return terminal, nil
}

func (r *TerminalRepository) GetBySerialNumber(serialNumber string) (domain.Terminal, error) {
    query := `
        SELECT id, serial_number, address, name
        FROM terminals
        WHERE serial_number = ?`

    var terminal domain.Terminal
    err := r.db.Get(&terminal, query, serialNumber)
    if err != nil {
        return domain.Terminal{}, fmt.Errorf("failed to get terminal by serial number: %w", err)
    }

    return terminal, nil
}

func (r *TerminalRepository) GetAll() ([]domain.Terminal, error) {
	query := `
        SELECT id, serial_number, address, name
        FROM terminals`

	var terminals []domain.Terminal
	err := r.db.Select(&terminals, query)
	if err != nil {
		return []domain.Terminal{}, fmt.Errorf("failed to get terminals: %w", err)
	}

	return terminals, nil
}

func (r *TerminalRepository) Update(id int, terminal domain.Terminal) error {
	var updates []string
	var args []interface{}

	if terminal.SerialNumber != "" {
		updates = append(updates, "serial_number = ?")
		args = append(args, terminal.SerialNumber)
	}
	if terminal.Address != "" {
		updates = append(updates, "address = ?")
		args = append(args, terminal.Address)
	}
	if terminal.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, terminal.Name)
	}

	if len(updates) == 0 {
		return fmt.Errorf("you must update at least 1 param")
	}

	query := fmt.Sprintf(`
        UPDATE terminals 
        SET %s 
        WHERE id = ?`, 
		strings.Join(updates, ", "))

	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update terminal: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("terminal with id %d not found", id)
	}

	return nil
}

func (r *TerminalRepository) Delete(id int) error {
	query := `
        DELETE FROM terminals
        WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete terminal: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("terminal with id %d not found", id)
	}

	return nil
}