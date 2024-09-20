package pg

import (
	"context"
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
