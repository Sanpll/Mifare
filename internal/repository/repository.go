package repository

import (
	"mifare/internal/domain"

	"github.com/jmoiron/sqlx"
)

type Key interface {
	Create(key domain.Key) (int, error)
	GetAll() ([]domain.Key, error)
	GetById(id int) (domain.Key, error)
	GetByValue(keyValue string) (domain.Key, error)
	Update(id int, key domain.Key) error
	Delete(id int) error
}

type Card interface {
	Create(card domain.Card) (int, error)
	GetById(id int) (domain.Card, error)
	GetByIdWithKey(id int) (domain.CardWithKey, error)
	GetByNumber(cardNumber string) (domain.CardWithKey, error)
	GetAll() ([]domain.Card, error)
	GetAllWithKeyValue() ([]domain.CardWithKey, error)
	Update(id int, card domain.Card, balanceProvided bool, isBlockedProvided bool) error
	Delete(id int) error
}

type Terminal interface {
	Create(terminal domain.Terminal) (int, error)
	GetById(id int) (domain.Terminal, error)
	GetBySerialNumber(serialNumber string) (domain.Terminal, error)
	GetAll() ([]domain.Terminal, error)
	Update(id int, terminal domain.Terminal) error
	Delete(id int) error
}

type Transaction interface {
	Create(transaction domain.Transaction) (int, error)
	GetById(id int) (domain.Transaction, error)
	GetAll() ([]domain.Transaction, error)
	Update(id int, user domain.Transaction) error
	Delete(id int) error
}

type User interface {
	Create(user domain.User) (int, error)
	GetByUsername(username string) (domain.User, error)
	GetById(id int) (domain.User, error)
	GetAll() ([]domain.User, error)
	Update(id int, user domain.User) error
	Delete(id int) error
}

type Repository struct {
	Key
	Card
	Terminal
	Transaction
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Key:         NewKeyRepository(db),
		Card:        NewCardRepository(db),
		Terminal:    NewTerminalRepository(db),
		Transaction: NewTransactionRepository(db),
		User:        NewUserRepository(db),
	}
}
