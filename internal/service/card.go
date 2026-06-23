package service

import (
	"database/sql"
    "errors"
	"fmt"
	"mifare/internal/domain"
	"mifare/internal/dto"
	"mifare/internal/repository"

	"github.com/shopspring/decimal"
)

type CardService struct {
	repo    repository.Card
	keyRepo repository.Key
}

func NewCardService(repo repository.Card, keyRepo repository.Key) *CardService {
	return &CardService{
		repo:    repo,
		keyRepo: keyRepo,
	}
}

var ErrCardNotFound = errors.New("card not found")

func (s *CardService) Create(input dto.CreateCardInput, username string) (int, error) {
	balance, err := decimal.NewFromString(input.Balance)
	if err != nil {
		return 0, fmt.Errorf("invalid balance format")
	}

	key, err := s.keyRepo.GetByValue(input.KeyValue)
	if err != nil {
		return 0, fmt.Errorf("failed to get key_id by value: %w", err)
	}

	card := domain.Card{
		CardNumber: input.CardNumber,
		Balance:    balance,
		IsBlocked:  input.IsBlocked,
		OwnerName:  username,
		KeyID:      key.ID,
	}

	return s.repo.Create(card)
}

func (s *CardService) GetAll() ([]dto.CardResponse, error) {
	domainCards, err := s.repo.GetAllWithKeyValue()
	if err != nil {
		return nil, fmt.Errorf("failed to get cards: %w", err)
	}

	dtoCards := make([]dto.CardResponse, 0, len(domainCards))
	for _, c := range domainCards {
		dtoCards = append(dtoCards, dto.CardResponse{
			ID:         c.ID,
			CardNumber: c.CardNumber,
			Balance:    c.Balance.String(),
			IsBlocked:  c.IsBlocked,
			OwnerName:  c.OwnerName,
			KeyValue:   c.KeyValue,
		})
	}

	return dtoCards, nil
}

func (s *CardService) GetById(id int) (dto.CardResponse, error) {
	card, err := s.repo.GetByIdWithKey(id)
	if err != nil {
		return dto.CardResponse{}, fmt.Errorf("failed to get card: %w", err)
	}

	return dto.CardResponse{
		ID:         card.ID,
		CardNumber: card.CardNumber,
		Balance:    card.Balance.String(),
		IsBlocked:  card.IsBlocked,
		OwnerName:  card.OwnerName,
		KeyValue:   card.KeyValue,
	}, nil
}

func (s *CardService) GetByNumber(cardNumber string) (dto.CardResponse, error) {
    card, err := s.repo.GetByNumber(cardNumber)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return dto.CardResponse{}, ErrCardNotFound
        }
        return dto.CardResponse{}, fmt.Errorf("failed to get card by number: %w", err)
    }
    return dto.CardResponse{
        ID:         card.ID,
        CardNumber: card.CardNumber,
        Balance:    card.Balance.String(),
        IsBlocked:  card.IsBlocked,
        OwnerName:  card.OwnerName,
        KeyValue:   card.KeyValue,
    }, nil
}

func (s *CardService) Update(id int, input dto.CardUpdate) error {
	balance, err := decimal.NewFromString(*input.Balance)
	if err != nil {
		return fmt.Errorf("invalid balance format")
	}

	card := domain.Card{}
	balanceProvided := false
	isBlockedProvided := false

	if input.Balance != nil {
		card.Balance = balance
		balanceProvided = true
	}
	if input.IsBlocked != nil {
		card.IsBlocked = *input.IsBlocked
		isBlockedProvided = true
	}
	if input.OwnerName != nil {
		card.OwnerName = *input.OwnerName
	}
	if input.KeyValue != nil {
		key, err := s.keyRepo.GetByValue(*input.KeyValue)
		if err != nil {
			return fmt.Errorf("failed to get key_id by value: %w", err)
		}
		card.KeyID = key.ID
	}

	return s.repo.Update(id, card, balanceProvided, isBlockedProvided)
}

func (s *CardService) Delete(id int) error {
	return s.repo.Delete(id)
}
