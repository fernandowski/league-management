package pg

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	domain "league-management/internal/domain/user"
	"log"
)

func Save(u *domain.User) error {
	pool := GetConnection()

	sql := "INSERT INTO league_management.users (email, password) VALUES ($1, $2)"

	_, err := pool.Exec(context.Background(), sql, u.Email, u.PasswordHash)
	if err != nil {
		log.Println(err)
	}

	return err
}

func FindById(emailAddress string) *domain.User {
	pool := GetConnection()

	sql := "SELECT email FROM league_management.users where email=$1"

	var email string

	err := pool.QueryRow(context.Background(), sql, emailAddress).Scan(&email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		panic(err)
	}

	return &domain.User{Email: email}
}
