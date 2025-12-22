package postgres

import (
	"context"
	"wimed/internal/application/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager struct {
	Pool *pgxpool.Pool
}

type txWrap struct {
	ctx context.Context
	tx pgx.Tx
}

func (t txWrap) Commit() error   { return t.tx.Commit(t.ctx) }
func (t txWrap) Rollback() error { return t.tx.Rollback(t.ctx) }

func (m TxManager) Begin(ctx context.Context) (ports.Tx, error) {
	tx, err := m.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return txWrap{ctx: ctx,tx: tx}, nil
}

func unwrapTx(tx ports.Tx) (pgx.Tx, bool) {
	w, ok := tx.(txWrap)
	if !ok {
		return nil, false
	}
	return w.tx, true
}