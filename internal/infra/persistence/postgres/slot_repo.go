package postgres

import (
	"context"
	"errors"
	"time"
	"wimed/internal/application/ports"
	"wimed/internal/domain/availabilityDomain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SlotRepository struct {
	Pool *pgxpool.Pool
}

func (r SlotRepository) GetByIDForUpdate(ctx context.Context, tx ports.Tx, id string) (*availabilityDomain.SlotDomain, error) {
	pgxTx, ok := unwrapTx(tx)
	if !ok {
		return nil, errors.New("invalid tx type")
	}

	const q = `
	SELECT id, doctor_id, start_at, end_at, status, created_at, updated_at
	FROM availability_slots
	WHERE id = $1
	FOR UPDATE
`
	var (
		slotID, DoctorID, status                 string
		startedAt, endedAt, createdAt, updatedAt time.Time
	)

	err := pgxTx.QueryRow(ctx, q, id).Scan(
		&slotID, &DoctorID, &startedAt, &endedAt, &status, &createdAt, &updatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("slot not found")
		}
		return nil, err
	}

	slot, err := availabilityDomain.RebuildSlotDomain(
		slotID,
		DoctorID,
		startedAt,
		endedAt,
		availabilityDomain.SlotStatus(status),
		createdAt,
		updatedAt,
	)

	if err != nil {
		return nil, err
	}
	return slot, nil
}

func (r SlotRepository) Update(ctx context.Context, tx ports.Tx, s *availabilityDomain.SlotDomain) error {
	pgxTx, ok := unwrapTx(tx)
	if !ok {
		return errors.New("invalid tx type")
	}

		const q = `
UPDATE availability_slots
SET status = $2,
    updated_at = $3
WHERE id = $1
`

	ct, err := pgxTx.Exec(ctx, q, s.ID(), string(s.Status()), s.UpdatedAt())
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("slot not fund")
	}
	return nil
}