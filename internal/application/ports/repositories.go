package ports

import (
	"context"
	"time"

	"wimed/internal/domain/appointmentDomain"
	"wimed/internal/domain/availabilityDomain"
	"wimed/internal/domain/paymentDomain"
)

type DoctorRepository interface {}

type PatientRepository interface {
	ExistsByID(ctx context.Context, tx Tx, id string) (bool, error)
}

type SlotLockRepository interface {
	GetByIDForUpdate(ctx context.Context, tx Tx, id string) (*availabilityDomain.SlotDomain, error)
	Update(ctx context.Context, tx Tx, s *availabilityDomain.SlotDomain) error
}

type SlotWriteRepoitory interface {
	CreateBatch(ctx context.Context, tx Tx, slots []*availabilityDomain.SlotDomain) error
}

type SlotReadRepository interface {
	ListAvailableByDoctor(ctx context.Context, tx Tx, doctorID string, from, to time.Time, limit int) ([]availabilityDomain.SlotDomain, error)
}

type AppointmentRepository interface {
	Create(ctx context.Context, tx Tx, a *appointmentDomain.AppointmentDomain) error
}

type PaymentRepository interface {
	Create(ctx context.Context, tx Tx, p *paymentDomain.PaymentDomain) error
}
