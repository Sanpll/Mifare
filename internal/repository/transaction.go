package repository

import (
	"fmt"
	"mifare/internal/domain"
	"strings"

	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) Create(transaction domain.Transaction) (int, error) {
	query := `
        INSERT INTO transactions (price, card_id, terminal_id, status)
        VALUES (?, ?, ?, ?)
        RETURNING id`

	var id int
	err := r.db.QueryRowx(query,
		transaction.Price,
		transaction.CardID,
		transaction.TerminalID,
		transaction.Status,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create transaction: %w", err)
	}

	return id, nil
}

func (r *TransactionRepository) GetById(id int) (domain.Transaction, error) {
	query := `
        SELECT id, price, card_id, terminal_id, status, created_at
        FROM transactions
        WHERE id = ?`

	var transaction domain.Transaction
	err := r.db.Get(&transaction, query, id)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

func (r *TransactionRepository) GetAll() ([]domain.Transaction, error) {
	query := `
        SELECT id, price, card_id, terminal_id, status, created_at
        FROM transactions`

	var transactions []domain.Transaction
	err := r.db.Select(&transactions, query)
	if err != nil {
		return []domain.Transaction{}, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}

func (r *TransactionRepository) Update(id int, transaction domain.Transaction) error {
	var updates []string
	var args []interface{}

	if transaction.CardID != 0 {
		updates = append(updates, "card_id = ?")
		args = append(args, transaction.CardID)
	}
	if transaction.Status != "" {
		updates = append(updates, "status = ?")
		args = append(args, transaction.Status)
	}
	if transaction.TerminalID != 0 {
		updates = append(updates, "terminal_id = ?")
		args = append(args, transaction.TerminalID)
	}

	if len(updates) == 0 {
		return fmt.Errorf("you must update at least 1 param")
	}

	query := fmt.Sprintf(`
        UPDATE transactions 
        SET %s 
        WHERE id = ?`, 
		strings.Join(updates, ", "))

	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("transaction with id %d not found", id)
	}

	return nil
}

func (r *TransactionRepository) Delete(id int) error {
	query := `
        DELETE FROM transactions
        WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("transaction with id %d not found", id)
	}

	return nil
}