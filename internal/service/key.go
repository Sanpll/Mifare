package service

import (
	"fmt"
	"mifare/internal/domain"
	"mifare/internal/dto"
	"mifare/internal/repository"
)

type KeyService struct {
	repo repository.Key
}

func NewKeyService(repo repository.Key) *KeyService {
	return &KeyService{
		repo: repo,
	}
}

func (s *KeyService) Create(input dto.CreateKeyInput) (int, error) {
	key := domain.Key{
		KeyValue:    input.KeyValue,
		KeyType:     input.KeyType,
		Description: input.Description,
	}

	return s.repo.Create(key)
}

func (s *KeyService) GetAll() ([]dto.KeyResponse, error) {
	domainKeys, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}

	dtoKeys := make([]dto.KeyResponse, 0, len(domainKeys))
	for _, k := range domainKeys {
		dtoKeys = append(dtoKeys, dto.KeyResponse{
			ID:          k.ID,
			KeyValue:    k.KeyValue,
			KeyType:     k.KeyType,
			Description: k.Description,
		})
	}

	return dtoKeys, nil
}

func (s *KeyService) GetById(id int) (dto.KeyResponse, error) {
	key, err := s.repo.GetById(id)
	if err != nil {
		return dto.KeyResponse{}, fmt.Errorf("failed to get key: %w", err)
	}

	return dto.KeyResponse{
		ID:          key.ID,
		KeyValue:    key.KeyValue,
		KeyType:     key.KeyType,
		Description: key.Description,
	}, nil
}

func (s *KeyService) GetByValue(keyValue string) (dto.KeyResponse, error) {
	if keyValue == "" {
		return dto.KeyResponse{}, fmt.Errorf("key_value is required")
	}

	key, err := s.repo.GetByValue(keyValue)
	if err != nil {
		return dto.KeyResponse{}, fmt.Errorf("failed to get key_id by value: %w", err)
	}

	return dto.KeyResponse{
		ID: key.ID,
		KeyValue: key.KeyValue,
		KeyType: key.KeyType,
		Description: key.Description,
	}, nil
}

func (s *KeyService) Update(id int, input dto.KeyUpdate) error {
	key := domain.Key{}

	if input.KeyValue != nil {
		key.KeyValue = *input.KeyValue
	}
	if input.KeyType != nil {
		key.KeyType = *input.KeyType
	}
	if input.Description != nil {
		key.Description = *input.Description
	}

	return s.repo.Update(id, key)
}

func (s *KeyService) Delete(id int) error {
	return s.repo.Delete(id)
}
