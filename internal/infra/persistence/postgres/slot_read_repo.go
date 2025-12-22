package postgres

import (
	"context"
	"errors"
	"time"

	"wimed/internal/application/ports"
	"wimed/internal/domain/availabilityDomain"
)

func (r SlotRepository) ListAvailableByDoctor(ctx context.Context, tx ports.Tx, doctorID string, from, to time.Time, limit int) ([]availabilityDomain.SlotDomain, error) {
	pgxTx, ok := unwrapTx(tx)
	if !ok {
		return nil, errors.New("invalid tx type")
	}
	if limit <= 0 {
		limit = 50
	}

	const q = `
SELECT id, doctor_id, start_at, end_at, status, created_at, updated_at
FROM availability_slots
WHERE doctor_id = $1
  AND status = 'AVAILABLE'
  AND start_at >= $2
  AND start_at <  $3
ORDER BY start_at ASC
LIMIT $4
`
	rows, err := pgxTx.Query(ctx, q, doctorID, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]availabilityDomain.SlotDomain, 0, limit)

	for rows.Next() {
		var (
			id, did, status string
			startedAt, endedAt, createdAt, updatedAt time.Time
		)
		if err := rows.Scan(&id, &did, &startedAt, &endedAt, &status, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		s, err := availabilityDomain.RebuildSlotDomain(
			id,
			did,
			startedAt,
			endedAt,
			availabilityDomain.SlotStatus(status),
			createdAt,
			updatedAt,
		)
		if err != nil {
			return nil, err
		}
		if s == nil {
			return nil, errors.New("failed to rebuild slot domain")
		}

		out = append(out, *s)
	}	
		
	if err := rows.Err(); err != nil {
			return nil, err
	}
		
	return out, nil
}