package pgsqltx

import (
	"context"
	"warehousesvc/internal/core/tx"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type pgsqlTx struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) tx.TransactionManager {
	return &pgsqlTx{db}
}

func (tx *pgsqlTx) Do(ctx context.Context, fn func(context.Context) error) (err error) {
	ctx, commit, err := tx.wrapContextAndGetCommitFn(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to wrap context and start transaction")
	}

	defer commit(&err)

	if err := fn(ctx); err != nil {
		return errors.Wrap(err, "failed to execute transactional fn")
	}

	return nil
}
