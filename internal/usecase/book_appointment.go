package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"wimed/internal/application/dto"
	"wimed/internal/application/ports"
	"wimed/internal/domain/appointment"
	"wimed/internal/domain/availability"
	"wimed/internal/domain/payment"
)

type BookAppointment struct {
	TxManager ports.TxManager

	Patients     ports.PatientRepository
	Slots        ports.SlotRepository
	Appointments ports.AppointmentRepository
	Payments     ports.PaymentRepository

	Now func() time.Time
}

func (uc *BookAppointment) Execute(ctx context.Context, in dto.BookAppointmentInput) (*dto.BookAppointmentOutput, error) {
	// validações de input (application-level)
	if strings.TrimSpace(in.AppointmentID) == "" {
		return nil, errors.New("appointment_id is required")
	}
	if strings.TrimSpace(in.PaymentID) == "" {
		return nil, errors.New("payment_id is required")
	}
	if strings.TrimSpace(in.SlotID) == "" {
		return nil, errors.New("slot_id is required")
	}
	if strings.TrimSpace(in.DoctorID) == "" {
		return nil, errors.New("doctor_id is required")
	}
	if strings.TrimSpace(in.PatientID) == "" {
		return nil, errors.New("patient_id is required")
	}
	if in.PriceCents < 0 {
		return nil, errors.New("price_cents must be >= 0")
	}

	now := time.Now()
	if uc.Now != nil {
		now = uc.Now()
	}

	tx, err := uc.TxManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	// padrão seguro de rollback
	defer func() { _ = tx.Rollback() }()

	// 1) garante que patient existe (no MVP, pode ser só "existe")
	ok, err := uc.Patients.ExistsByID(ctx, tx, in.PatientID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("patient not found")
	}

	// 2) carrega slot com lock (FOR UPDATE)
	slot, err := uc.Slots.GetByIDForUpdate(ctx, tx, in.SlotID)
	if err != nil {
		return nil, err
	}

	// 3) checa invariantes de negócio aqui
	if slot.DoctorID() != in.DoctorID {
		return nil, errors.New("slot does not belong to doctor")
	}
	if slot.Status() != availability.SlotAvailable {
		return nil, errors.New("slot is not available")
	}

	// 4) marca slot como booked (regra no domínio)
	if err := slot.MarkBooked(now); err != nil {
		return nil, err
	}
	if err := uc.Slots.Update(ctx, tx, slot); err != nil {
		return nil, err
	}

	// 5) cria appointment (congela preço)
	a, err := appointment.NewAppointment(
		in.AppointmentID,
		in.DoctorID,
		in.PatientID,
		in.SlotID,
		in.PriceCents,
		appointment.StatusScheduled,
		now,
	)
	if err != nil {
		return nil, err
	}
	if err := uc.Appointments.Create(ctx, tx, a); err != nil {
		return nil, err
	}

	// 6) cria payment PENDING
	provider, err := parseProvider(in.PaymentProvider)
	if err != nil {
		return nil, err
	}

	p, err := payment.NewPayment(
		in.PaymentID,
		in.AppointmentID,
		provider,
		in.PriceCents,
		payment.StatusPending,
		in.ExternalRef,
		now,
	)
	if err != nil {
		return nil, err
	}
	if err := uc.Payments.Create(ctx, tx, p); err != nil {
		return nil, err
	}

	// 7) commit transação
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.BookAppointmentOutput{
		AppointmentID: a.ID(),
		PaymentID:     p.ID(),
		Status:        string(a.Status()),
	}, nil
}

func parseProvider(raw string) (payment.Provider, error) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "STRIPE":
		return payment.ProviderStripe, nil
	case "MERCADOPAGO":
		return payment.ProviderMercadoPago, nil
	case "MANUAL", "":
		return payment.ProviderManual, nil
	default:
		return "", errors.New("invalid payment provider")
	}
}
