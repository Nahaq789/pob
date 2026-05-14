package port

import (
	"context"
	"pob/user/internal/shared"
)

type TransactionManager struct {
	db *shared.DBClient
}

func NewTransactionManager(d *shared.DBClient) *TransactionManager {
	return &TransactionManager{db: d}
}

func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.GetClient().Begin(ctx)
	if err != nil {
		return err
	}

	txCtx := shared.TxWithContext(ctx, tx)

	if err := fn(txCtx); err != nil {
		tx.Rollback(txCtx)
		return err
	}
	return tx.Commit(txCtx)
}
