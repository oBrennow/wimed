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
	tx pgx.Tx
}

func (t txWrap) Commit() error   { return t.tx.Commit(context.Background()) }
func (t txWrap) Rollback() error { return t.tx.Rollback(context.Background()) }

func (m TxManager) Begin(ctx context.Context) (ports.Tx, error) {
	tx, err := m.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return txWrap{tx: tx}, nil
}
