package service

import (
	"mifare/internal/dto"
	"mifare/internal/repository"
)

type Authorization interface {
	GenerateToken(input dto.SignInInput) (string, error)
	ParseToken(token string) (int, string, bool, error)
}

type Key interface {
	Create(input dto.CreateKeyInput) (int, error)
	GetAll() ([]dto.KeyResponse, error)
	GetById(id int) (dto.KeyResponse, error)
	GetByValue(keyValue string) (dto.KeyResponse, error)
	Update(id int, input dto.KeyUpdate) error
	Delete(id int) error
}

type Card interface {
	Create(input dto.CreateCardInput, username string) (int, error)
	GetAll() ([]dto.CardResponse, error)
	GetById(id int) (dto.CardResponse, error)
	GetByNumber(cardNumber string) (dto.CardResponse, error)
	Update(id int, input dto.CardUpdate) error
	Delete(id int) error
}

type Terminal interface {
	Create(input dto.CreateTerminalInput) (int, error)
	GetAll() ([]dto.TerminalResponse, error)
	GetById(id int) (dto.TerminalResponse, error)
	Update(id int, input dto.TerminalUpdate) error
	Delete(id int) error
}

type Transaction interface {
	Authorize(input dto.AuthorizeTransactionInput) (dto.AuthorizeTransactionResponse, error)
	Create(input dto.CreateTransactionInput) (int, error)
	GetAll() ([]dto.TransactionResponse, error)
	GetById(id int) (dto.TransactionResponse, error)
	Update(id int, input dto.TransactionUpdate) error
	Delete(id int) error
}

type User interface {
	Create(input dto.SignUpInput) (int, error)
	GetAll() ([]dto.UserResponse, error)
	GetById(id int) (dto.UserResponse, error)
	Update(id int, input dto.UserUpdate) error
	Delete(id int) error
}

type Service struct {
	Authorization
	Key
	Card
	Terminal
	Transaction
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.User),
		Key: NewKeyService(repo.Key),
		Card: NewCardService(repo.Card, repo.Key),
		Terminal: NewTerminalService(repo.Terminal),
		Transaction: NewTransactionService(repo.Transaction, repo.Card, repo.Terminal),
		User: NewUserService(repo.User),
	}
}
