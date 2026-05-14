package shared

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type contextTransactionKey struct{}

func TxWithContext(ctx context.Context, tx pgx.Tx) context.Context {
	ctx = context.WithValue(ctx, contextTransactionKey{}, tx)
	return ctx
}

func TxFromContext(ctx context.Context) pgx.Tx {
	tx, ok := ctx.Value(contextTransactionKey{}).(pgx.Tx)
	if !ok {
		return nil
	}
	return tx
}
