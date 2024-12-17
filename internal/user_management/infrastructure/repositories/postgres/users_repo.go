package pg

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"league-management/internal/shared/database"
	"league-management/internal/user_management/domain/user"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) Save(u *domain.User) error {
	pool := database.GetConnection()

	sql := "INSERT INTO league_management.users (email, password) VALUES ($1, $2)"

	_, err := pool.Exec(context.Background(), sql, u.Email, u.PasswordHash)
	if err != nil {
	}

	return err
}

func (ur *UserRepository) FindByEmail(emailAddress string) *domain.User {
	pool := database.GetConnection()

	sql := "SELECT id, email, password FROM league_management.users where email=$1"

	var email string
	var password string
	var id string

	err := pool.QueryRow(context.Background(), sql, emailAddress).Scan(&id, &email, &password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		panic(err)
	}

	return &domain.User{Id: id, Email: email, PasswordHash: password}
}

func (ur *UserRepository) FindById(userId string) (*domain.User, error) {
	var pool = database.GetConnection()

	var sql = "SELECT id, email FROM league_management.users where id=$1"

	var email string
	var id string

	err := pool.QueryRow(context.Background(), sql, userId).Scan(&id, &email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		panic(err)
	}

	return &domain.User{Id: id, Email: email}, nil
}
