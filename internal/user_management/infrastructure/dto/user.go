package dto

import (
	"league-management/internal/user_management/domain/user"
)

type UserDto struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func FromDomain(user *domain.User) *UserDto {
	return &UserDto{
		Id:    user.Id,
		Email: user.Email,
	}
}
