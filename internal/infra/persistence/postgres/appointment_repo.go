package postgres

import (
	"context"
	"errors"

	"wimed/internal/application/ports"
	"wimed/internal/domain/appointmentDomain"

	"github.com/jackc/pgx/v5/pgconn"
)

type AppointmentRepository struct{}

func (r AppointmentRepository) Create(ctx context.Context, tx ports.Tx, a *appointmentDomain.AppointmentDomain) error {
	pgxTx, ok := unwrapTx(tx)
	if !ok {
		return errors.New("invalid tx type")
	}

	const q = `
INSERT INTO appointments (id, doctor_id, patient_id, slot_id, price_cents, status, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
`
	_, err := pgxTx.Exec(ctx, q,
	a.ID(),
	a.DoctorID(),
	a.PatientID(),
	a.SlotID(),
	a.PriceCents(),
	string(a.Status()),
	a.CreatedAt(),
	a.UpdatedAt(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return errors.New("appointment already exists for this slot")
			}
		}
		return err
	}
	return nil
}