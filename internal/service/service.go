package service

import (
	"mifare/internal/domain"
	"mifare/internal/repository"
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

type Service struct {
	Authorization
	Card
	Key
	Terminal
	Transaction
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}