package service

import (
	"fmt"
	"mifare/internal/domain"
	"mifare/internal/dto"
	"mifare/internal/repository"
)

type TerminalService struct {
	repo repository.Terminal
}

func NewTerminalService(repo repository.Terminal) *TerminalService {
	return &TerminalService{
		repo: repo,
	}
}

func (s *TerminalService) Create(input dto.CreateTerminalInput) (int, error) {
	terminal := domain.Terminal{
		SerialNumber: input.SerialNumber,
		Address:      input.Address,
		Name:         input.Name,
	}

	return s.repo.Create(terminal)
}

func (s *TerminalService) GetAll() ([]dto.TerminalResponse, error) {
	domainTerminals, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get terminals: %w", err)
	}

	dtoTerminals := make([]dto.TerminalResponse, 0, len(domainTerminals))
	for _, t := range domainTerminals {
		dtoTerminals = append(dtoTerminals, dto.TerminalResponse{
			ID:           t.ID,
			SerialNumber: t.SerialNumber,
			Address:      t.Address,
			Name:         t.Name,
		})
	}

	return dtoTerminals, nil
}

func (s *TerminalService) GetById(id int) (dto.TerminalResponse, error) {
	terminal, err := s.repo.GetById(id)
	if err != nil {
		return dto.TerminalResponse{}, fmt.Errorf("failed to get terminal: %w", err)
	}

	return dto.TerminalResponse{
		ID:           terminal.ID,
		SerialNumber: terminal.SerialNumber,
		Address:      terminal.Address,
		Name:         terminal.Name,
	}, nil
}

func (s *TerminalService) Update(id int, input dto.TerminalUpdate) error {
	terminal := domain.Terminal{}

	if input.SerialNumber != nil {
		terminal.SerialNumber = *input.SerialNumber
	}
	if input.Address != nil {
		terminal.Address = *input.Address
	}
	if input.Name != nil {
		terminal.Name = *input.Name
	}

	return s.repo.Update(id, terminal)
}

func (s *TerminalService) Delete(id int) error {
	return s.repo.Delete(id)
}