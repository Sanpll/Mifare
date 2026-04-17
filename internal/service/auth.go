package service

import (
	"errors"
	"fmt"
	"mifare/internal/dto"
	"mifare/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const(
	tokenTTL = 2 * time.Hour
	signingKey = "super-secret-key-which-requires-at-least-32-characters-inside"
)

type AuthService struct {
	repo repository.User
}

func NewAuthService(repo repository.User) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) GenerateToken(input dto.SignInInput) (string, error) {
	
	user, err := s.repo.GetByUsername(input.Username)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":      "mifare-server",
		"sub":      fmt.Sprintf("%d", user.ID),
		"user_id":  user.ID,
		"username": user.Username,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(tokenTTL).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(inputToken string) (int, string, bool, error) {
	token, err := jwt.Parse(inputToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, "", false, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return 0, "", false, errors.New("token is expired or invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", false, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", false, errors.New("user_id not found in token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return 0, "", false, errors.New("username not found in token")
	}

	isAdmin, ok := claims["is_admin"].(bool)
	if !ok {
		isAdmin = false
	}

	return int(userID), username, isAdmin, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
