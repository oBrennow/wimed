package ports

import "context"

type TX interface {
	Commit() error
	Rollback() error
}

type TxManager interface {
	Begin(ctx context.Context) (TX, error)
}
