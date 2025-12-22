package postgres

import (
	"context"
	"errors"
	
	"wimed/internal/application/ports"
)

type PatientRepository struct{}

func (r PatientRepository) ExistsByID(ctx context.Context, tx ports.Tx, id string) (bool, error) {
	pgxTx, ok := unwrapTx(tx)
	if !ok {
		return false, errors.New("invalid tx type")
	}

	const q = `SELECT EXISTS(SELECT 1 FROM patients WHERE id = $1)`
	var exists bool
	if err := pgxTx.QueryRow(ctx, q, id).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}