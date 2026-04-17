package repository

import (
	"fmt"
	"mifare/internal/domain"
	"strings"

	"github.com/jmoiron/sqlx"
)

type CardRepository struct {
	db *sqlx.DB
}

func NewCardRepository(db *sqlx.DB) *CardRepository {
	return &CardRepository{
		db: db,
	}
}

func (r *CardRepository) Create(card domain.Card) (int, error) {
	query := `
        INSERT INTO cards (card_number, balance, is_blocked, owner_name, key_id)
        VALUES (?, ?, ?, ?, ?)
        RETURNING id`

	var id int
	err := r.db.QueryRowx(query,
		card.CardNumber,
		card.Balance,
		card.IsBlocked,
		card.OwnerName,
		card.KeyID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create card: %w", err)
	}

	return id, nil
}

func (r *CardRepository) GetById(id int) (domain.Card, error) {
	query := `
        SELECT *
        FROM cards
        WHERE id = ?`

	var card domain.Card
	err := r.db.Get(&card, query, id)
	if err != nil {
		return domain.Card{}, fmt.Errorf("failed to get card: %w", err)
	}

	return card, nil
}

func (r *CardRepository) GetByIdWithKey(id int) (domain.CardWithKey, error) {
	query := `
        SELECT 
            c.id, c.card_number, c.balance, c.is_blocked, c.owner_name, c.key_id,
            k.key_value
        FROM cards c
        LEFT JOIN keys k ON c.key_id = k.id
		WHERE c.id = ?
        ORDER BY c.id`

	var card domain.CardWithKey
	err := r.db.Get(&card, query, id)
	if err != nil {
		return domain.CardWithKey{}, fmt.Errorf("failed to get card: %w", err)
	}

	return card, nil
}

func (r *CardRepository) GetByNumber(cardNumber string) (domain.CardWithKey, error) {
    query := `
        SELECT 
            c.id, c.card_number, c.balance, c.is_blocked, c.owner_name, c.key_id,
            k.key_value
        FROM cards c
        LEFT JOIN keys k ON c.key_id = k.id
        WHERE c.card_number = ?`

    var card domain.CardWithKey
    err := r.db.Get(&card, query, cardNumber)
    if err != nil {
        return domain.CardWithKey{}, fmt.Errorf("failed to get card by number: %w", err)
    }

    return card, nil
}

func (r *CardRepository) GetAll() ([]domain.Card, error) {
	query := `
        SELECT *
        FROM cards`

	var cards []domain.Card
	err := r.db.Select(&cards, query)
	if err != nil {
		return []domain.Card{}, fmt.Errorf("failed to get cards: %w", err)
	}

	return cards, nil
}

func (r *CardRepository) GetAllWithKeyValue() ([]domain.CardWithKey, error) {
	query := `
        SELECT 
            c.id, c.card_number, c.balance, c.is_blocked, c.owner_name, c.key_id,
            k.key_value
        FROM cards c
        LEFT JOIN keys k ON c.key_id = k.id
        ORDER BY c.id`

	var cards []domain.CardWithKey
	err := r.db.Select(&cards, query)
	if err != nil {
		return []domain.CardWithKey{}, fmt.Errorf("failed to get cards with key_value: %w", err)
	}

	return cards, nil
}

func (r *CardRepository) Update(id int, card domain.Card, balanceProvided bool, isBlockedProvided bool) error {
	var updates []string
	var args []interface{}

	if balanceProvided || !card.Balance.IsZero() {
		updates = append(updates, "balance = ?")
		args = append(args, card.Balance)
	}
	if isBlockedProvided {
		updates = append(updates, "is_blocked = ?")
		args = append(args, card.IsBlocked)
	}
	if card.OwnerName != "" {
		updates = append(updates, "owner_name = ?")
		args = append(args, card.OwnerName)
	}
	if card.KeyID != 0 {
		updates = append(updates, "key_id = ?")
		args = append(args, card.KeyID)
	}

	if len(updates) == 0 {
		return fmt.Errorf("you must update at least 1 param")
	}

	query := fmt.Sprintf(`
        UPDATE cards 
        SET %s 
        WHERE id = ?`, 
		strings.Join(updates, ", "))

	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update card: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("card with id %d not found", id)
	}

	return nil
}

func (r *CardRepository) Delete(id int) error {
	query := `
        DELETE FROM cards
        WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("card with id %d not found", id)
	}

	return nil
}