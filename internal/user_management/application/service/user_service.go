package service

import (
	"errors"
	"github.com/google/uuid"
	"league-management/internal/user_management/domain/user"
	"league-management/internal/user_management/infrastructure/crypto"
	"league-management/internal/user_management/infrastructure/repositories/postgres"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

var userRepository = pg.NewUserRepository()

func (us *UserService) RegisterUser(email string, password string) (*domain.User, error) {

	var user = userRepository.FindByEmail(email)

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

	err := userRepository.Save(newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil

}

func (us *UserService) Login(email string, password string) (string, error) {
	var existingUser = userRepository.FindByEmail(email)

	if existingUser == nil {
		return "", errors.New("invalid user name or password")
	}

	if !crypto.VerifyPassword(password, existingUser.PasswordHash) {
		return "", errors.New("invalid user name or password")
	}

	service := NewAuthService()
	signedToken, err := service.GenerateJWT(existingUser)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
