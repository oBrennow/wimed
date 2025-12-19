package postgres


import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager struct {
	Pool	*pgxpool.Pool
}



