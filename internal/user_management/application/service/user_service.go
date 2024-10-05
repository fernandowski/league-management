package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"league-management/internal/user_management/internal/domain/user"
	"league-management/internal/user_management/internal/infrastructure/crypto"
	"league-management/internal/user_management/internal/infrastructure/repositories/postgres"
	"time"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) RegisterUser(email string, password string) (*domain.User, error) {

	var user = pg.FindByEmail(email)

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

func (us *UserService) Login(email string, password string) (string, error) {
	var existingUser = pg.FindByEmail(email)

	if existingUser == nil {
		return "", errors.New("invalid user name or password")
	}

	if !crypto.VerifyPassword(password, existingUser.PasswordHash) {
		return "", errors.New("invalid user name or password")
	}

	var jwtKey = []byte("my_secret_key")

	var claims = jwt.MapClaims{
		"authorized": true,
		"user_email": existingUser.Email,
		"exp":        time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
