package domain

import (
	"errors"
	"strings"
)

type User struct {
	Id                 string
	Email              string
	ContactInformation ContactInformation
	PasswordHash       string
}

func NewUser(id string, email string, passwordHash string) (*User, error) {

	if strings.TrimSpace(email) == "" {
		return &User{Email: email}, errors.New("email value cannot be empty")
	}

	if strings.TrimSpace(passwordHash) == "" {
		return &User{Email: email}, errors.New("password value cannot be empty")
	}

	if strings.TrimSpace(id) == "" {
		return &User{Email: email}, errors.New("password value cannot be empty")
	}

	return &User{
		Id:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}
