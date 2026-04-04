package database

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := GetConnection().BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
