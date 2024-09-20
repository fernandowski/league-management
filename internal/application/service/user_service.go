package service

import (
	"github.com/google/uuid"
	domain "league-management/internal/domain/user"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) RegisterUser(email string, password string) (*domain.User, error) {
	newUserId := uuid.New().String()
	newUser, err := domain.NewUser(newUserId, email, password)

	if err != nil {
		return nil, err
	}

	return newUser, nil

}
