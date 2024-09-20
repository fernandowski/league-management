package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"sync"
)

var (
	dbPool     *pgxpool.Pool
	dbPoolOnce sync.Once
)

func GetConnection() *pgxpool.Pool {
	dbPoolOnce.Do(func() {
		var err error
		dbPool, err = pgxpool.New(context.Background(), "postgres://root:root@db:5432/league?sslmode=disable")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
			os.Exit(1)
		}
	})

	return dbPool
}
