package postgres

import (
	"context"
	"errors"
	
	"wimed/internal/application/ports"
	"wimed/internal/domain/paymentDomain"

	"github.com/jackc/pgx/v5/pgconn"
)

type PaymentRepository struct {}

func (r PaymentRepository) Create(ctx context.Context, tx ports.Tx, p *paymentDomain.PaymentDomain) error {
	pgxTx, ok := unwrapTx(tx)
	if !ok {
		return errors.New("invalid tx type")
	}

		const q = `
INSERT INTO payments (id, appointment_id, provider, amount_cents, status, external_ref, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
`
	_, err := pgxTx.Exec(ctx, q,
	p.ID(),
	p.AppointmentID(),
	string(p.Provider()),
	p.AmountCents(),
	string(p.Status()),
	p.ExternalRef(),
	p.CreatedAt(),
	p.UpdatedAt(),
	)
	if err != nil {
		// payments.appointment_id é UNIQUE no seu DDL.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return errors.New("payment already exists for this appointment")
			}
		}
		return err
	}

	return nil
}