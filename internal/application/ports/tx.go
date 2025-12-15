package ports

import "context"


type tx interface {
	Commit() error
	Rollback() error
}

type TxManager interface {
	Begin (ctx context.Context) (Tx, error)
}