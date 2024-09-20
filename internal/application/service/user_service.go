package service

import (
	"github.com/google/uuid"
	domain "league-management/internal/domain/user"
	"league-management/internal/infrastructure/crypto"
	pg "league-management/internal/infrastructure/repositories/postgres"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) RegisterUser(email string, password string) (*domain.User, error) {
	newUserId := uuid.New().String()

	newUser, userError := domain.NewUser(newUserId, email, password)

	if userError != nil {
		return nil, userError
	}

	hash, passwordHashErr := crypto.HashPassword(password)

	if passwordHashErr != nil {
		return nil, passwordHashErr
	}

	newUser.PasswordHash = hash

	err := pg.Save(newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil

}
