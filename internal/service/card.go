package service

import (
	"fmt"
	"mifare/internal/domain"
	"mifare/internal/dto"
	"mifare/internal/repository"
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

func (s *CardService) Create(input dto.CreateCardInput, username string) (int, error) {
	key, err := s.keyRepo.GetByValue(input.KeyValue)
	if err != nil {
		return 0, fmt.Errorf("failed to get key_id by value: %w", err)
	}

	card := domain.Card{
		CardNumber: input.CardNumber,
		Balance:    input.Balance,
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
			Balance:    c.Balance,
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
		Balance:    card.Balance,
		IsBlocked:  card.IsBlocked,
		OwnerName:  card.OwnerName,
		KeyValue:   card.KeyValue,
	}, nil
}

func (s *CardService) Update(id int, input dto.CardUpdate) error {
	card := domain.Card{}
	balanceProvided := false
	isBlockedProvided := false

	if input.Balance != nil {
		card.Balance = *input.Balance
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