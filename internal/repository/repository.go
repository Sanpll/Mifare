package repository

import (
	"mifare/internal/domain"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
}

type Card interface {
}

type Key interface {
}

type Terminal interface {
}

type Transaction interface {
}

type User interface {
}

type Repository struct {
	Authorization
	Card
	Key
	Terminal
	Transaction
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSQLite(db),
	}
}