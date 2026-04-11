package service

import (
	"mifare/internal/domain"
	"mifare/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(user domain.User) (int, error) {
	var err error
	user.PasswordHash, err = hashPassword(user.PasswordHash)
	if err != nil {
		return -1, err
	}
	
	return s.repo.CreateUser(user)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}