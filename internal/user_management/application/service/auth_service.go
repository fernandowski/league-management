package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	domain "league-management/internal/user_management/domain/user"
	"time"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (as *AuthService) GenerateJWT(user *domain.User) (string, error) {
	var jwtKey = []byte("my_secret_key")

	var claims = jwt.MapClaims{
		"authorized": true,
		"user_email": user.Email,
		"exp":        time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (as *AuthService) ValidateJWT(signedJWT string) (jwt.MapClaims, error) {
	var jwtKey = []byte("my_secret_key")

	token, err := jwt.Parse(signedJWT, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("invalid JWT")
		}

		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("error Parsing JWT")
}
