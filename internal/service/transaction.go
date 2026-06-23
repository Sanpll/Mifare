package service

import (
	"fmt"
	"mifare/internal/domain"
	"mifare/internal/dto"
	"mifare/internal/repository"

	"github.com/shopspring/decimal"
)

type TransactionService struct {
	repo         repository.Transaction
	cardRepo     repository.Card
	terminalRepo repository.Terminal
}

func NewTransactionService(repo repository.Transaction, cardRepo repository.Card, terminalRepo repository.Terminal) *TransactionService {
	return &TransactionService{
		repo:         repo,
		cardRepo:     cardRepo,
		terminalRepo: terminalRepo,
	}
}

func (s *TransactionService) Authorize(input dto.AuthorizeTransactionInput) (dto.AuthorizeTransactionResponse, error) {
	price, err := decimal.NewFromString(input.Price)
	if err != nil {
		return dto.AuthorizeTransactionResponse{
			Authorized: false,
			Message: "incorrect price format",
		}, nil
	}

	if price.LessThan(decimal.Zero) {
		return dto.AuthorizeTransactionResponse{
			Authorized: false,
			Message: "incorrect price value",
		}, nil
	}

	card, err := s.cardRepo.GetByNumber(input.CardNumber)
	if err != nil {
		return dto.AuthorizeTransactionResponse{
			Authorized: false,
			Message: "card not found",
		}, nil
	}

	if card.IsBlocked {
		return dto.AuthorizeTransactionResponse{
			Authorized: false,
			Message: "card is blocked",
		}, nil
	}

	if card.Balance.LessThan(price) {
		return dto.AuthorizeTransactionResponse{
			Authorized: false,
			Message: "insufficient balance",
		}, nil
	}

	if input.TerminalSerialNumber != "" {
		_, err := s.terminalRepo.GetBySerialNumber(input.TerminalSerialNumber)
		if err != nil {
			return dto.AuthorizeTransactionResponse{
				Authorized: false,
				Message: "terminal not found",
			}, nil
		}
	}

	return dto.AuthorizeTransactionResponse{
		Authorized: true,
		Message: "authorized",
	}, nil
}

func (s *TransactionService) Create(input dto.CreateTransactionInput) (int, error) {
	price, err := decimal.NewFromString(input.Price)
	if err != nil {
		return 0, fmt.Errorf("invalid price format")
	}

	card, err := s.cardRepo.GetByNumber(input.CardNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to get card by card number: %w", err)
	}

	terminal, err := s.terminalRepo.GetBySerialNumber(input.TerminalSerialNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to get terminal by serial number: %w", err)
	}

	transaction := domain.Transaction{
		Price:      price,
		CardID:     card.ID,
		TerminalID: terminal.ID,
		Status:     input.Status,
	}

	return s.repo.Create(transaction)
}

func (s *TransactionService) GetAll() ([]dto.TransactionResponse, error) {
	domainTransactions, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	dtoTransactions := make([]dto.TransactionResponse, 0, len(domainTransactions))
	for _, t := range domainTransactions {
		card, err := s.cardRepo.GetById(int(t.CardID))
		if err != nil {
			return []dto.TransactionResponse{}, fmt.Errorf("failed to get card by id: %w", err)
		}

		terminal, err := s.terminalRepo.GetById(int(t.TerminalID))
		if err != nil {
			return []dto.TransactionResponse{}, fmt.Errorf("failed to get terminal by id: %w", err)
		}

		dtoTransactions = append(dtoTransactions, dto.TransactionResponse{
			ID:                   t.ID,
			CardNumber:           card.CardNumber,
			Price:                t.Price.String(),
			TerminalSerialNumber: terminal.SerialNumber,
		})
	}

	return dtoTransactions, nil
}

func (s *TransactionService) GetById(id int) (dto.TransactionResponse, error) {
	transaction, err := s.repo.GetById(id)
	if err != nil {
		return dto.TransactionResponse{}, fmt.Errorf("failed to get transaction: %w", err)
	}

	card, err := s.cardRepo.GetById(int(transaction.CardID))
	if err != nil {
		return dto.TransactionResponse{}, fmt.Errorf("failed to get card by id: %w", err)
	}

	terminal, err := s.terminalRepo.GetById(int(transaction.TerminalID))
	if err != nil {
		return dto.TransactionResponse{}, fmt.Errorf("failed to get terminal by id: %w", err)
	}

	return dto.TransactionResponse{
		ID:                   transaction.ID,
		CardNumber:           card.CardNumber,
		Price:                transaction.Price.String(),
		TerminalSerialNumber: terminal.SerialNumber,
	}, nil
}

func (s *TransactionService) Update(id int, input dto.TransactionUpdate) error {
	transaction := domain.Transaction{}

	if input.CardNumber != nil {
		card, err := s.cardRepo.GetByNumber(*input.CardNumber)
		if err != nil {
			return fmt.Errorf("failed to get card by card number: %w", err)
		}

		transaction.CardID = card.ID
	}
	if input.Status != nil {
		transaction.Status = *input.Status
	}
	if input.TerminalSerialNumber != nil {
		terminal, err := s.terminalRepo.GetBySerialNumber(*input.TerminalSerialNumber)
		if err != nil {
			return fmt.Errorf("failed to get terminal by serial number: %w", err)
		}

		transaction.TerminalID = terminal.ID
	}

	return s.repo.Update(id, transaction)
}

func (s *TransactionService) Delete(id int) error {
	return s.repo.Delete(id)
}
