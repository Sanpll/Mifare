package service

import (
	"fmt"
	"mifare/internal/domain"
	"mifare/internal/dto"
	"mifare/internal/repository"
	"strings"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(input dto.SignUpInput) (int, error) {
	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	user := domain.User{
		Username:     input.Username,
		PasswordHash: hashedPassword,
		IsAdmin:      input.IsAdmin,
	}

	id, err := s.repo.Create(user)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") ||
			strings.Contains(err.Error(), "users.username") ||
			strings.Contains(err.Error(), "constraint failed") {
			return 0, fmt.Errorf("username '%s' is already taken", input.Username)
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (s *UserService) GetAll() ([]dto.UserResponse, error) {
	domainUsers, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	dtoUsers := make([]dto.UserResponse, 0, len(domainUsers))
	for _, u := range domainUsers {
		dtoUsers = append(dtoUsers, dto.UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			IsAdmin:   u.IsAdmin,
			CreatedAt: u.CreatedAt,
		})
	}

	return dtoUsers, nil
}

func (s *UserService) GetById(id int) (dto.UserResponse, error) {
	user, err := s.repo.GetById(id)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserService) Update(id int, input dto.UserUpdate) error {
	if input.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	user := domain.User{
		Username: input.Username,
	}

	return s.repo.Update(id, user)
}

func (s *UserService) Delete(id int) error {
	return s.repo.Delete(id)
}
