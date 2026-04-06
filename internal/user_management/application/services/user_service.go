package services

import (
	"errors"
	"github.com/google/uuid"
	"league-management/internal/user_management/domain"
	"league-management/internal/user_management/infrastructure/crypto"
)

type UserService struct {
	userRepository userRepository
	authService    authService
}

type userRepository interface {
	Save(*domain.User) error
	FindByEmail(string) *domain.User
}

type authService interface {
	GenerateJWT(*domain.User) (string, error)
}

func NewUserService(userRepository userRepository, authService authService) *UserService {
	return &UserService{
		userRepository: userRepository,
		authService:    authService,
	}
}

func (us *UserService) RegisterUser(email string, password string) (*domain.User, error) {

	var user = us.userRepository.FindByEmail(email)

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

	err := us.userRepository.Save(newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil

}

func (us *UserService) Login(email string, password string) (string, error) {
	var existingUser = us.userRepository.FindByEmail(email)

	if existingUser == nil {
		return "", errors.New("invalid user name or password")
	}

	if !crypto.VerifyPassword(password, existingUser.PasswordHash) {
		return "", errors.New("invalid user name or password")
	}

	signedToken, err := us.authService.GenerateJWT(existingUser)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
