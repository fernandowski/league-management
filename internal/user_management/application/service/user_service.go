package service

import (
	"errors"
	"github.com/google/uuid"
	"league-management/internal/user_management/internal/domain/user"
	"league-management/internal/user_management/internal/infrastructure/crypto"
	"league-management/internal/user_management/internal/infrastructure/repositories/postgres"
	"log"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) RegisterUser(email string, password string) (*domain.User, error) {

	var user = pg.FindById(email)
	log.Print(user)
	if user != nil {
		return nil, errors.New("email is taken")
	}

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
