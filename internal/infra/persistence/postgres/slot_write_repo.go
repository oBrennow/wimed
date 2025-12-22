package postgres

import (
	"context"
	"errors"
	"wimed/internal/application/ports"
	"wimed/internal/domain/availabilityDomain"

	"github.com/jackc/pgx/v5"
)

func (r SlotRepository) CreateBatch (ctx context.Context, tx ports.Tx, slots []*availabilityDomain.SlotDomain) error {
	pgxTx, ok := unwrapTx(tx)
	if !ok {
		return errors.New("invalid tx type")
	}
	if len(slots) == 0 {
		return nil
	}

	const q = `
INSERT INTO availability_slots (id, doctor_id, start_at, end_at, status, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7)
ON CONFLICT (id) DO NOTHING
`

	batch := &pgx.Batch{}
	for _, s := range slots {
		batch.Queue(q, 
		s.ID(),
		s.DoctorID(),
		s.StartedAt(),
		s.EndedAt(),
		string(s.Status()),
		s.CreatedAt(),
		s.UpdatedAt(),
		)
	}

	br := pgxTx.SendBatch(ctx, batch)
	defer br.Close()

	for range slots {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}