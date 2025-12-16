package usecase

import (
	"context"
	"errors"
	"strings"
	"time"
	"wimed/internal/application/dto"
	"wimed/internal/domain/appointmentDomain"
	"wimed/internal/domain/availabilityDomain"
	"wimed/internal/domain/paymentDomain"

	"wimed/internal/application/ports"
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
	if strings.TrimSpace(in.AppointmentID) == "" {
		return nil, errors.New("appointment_id is required")
	}
	if strings.TrimSpace(in.PatientID) == "" {
		return nil, errors.New("patient_id is required")
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
		return nil, errors.New("price-cents must be positive")
	}

	now := time.Now()
	if uc.Now != nil {
		now = uc.Now()
	}

	tx, err := uc.TxManager.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() { _ = tx.Rollback() }()

	ok, err := uc.Patients.ExistsByID(ctx, tx, in.PatientID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("patient not found")
	}

	slot, err := uc.Slots.GetIDForUpate(ctx, tx, in.SlotID)
	if err != nil {
		return nil, err
	}

	if slot.DoctorID() != in.DoctorID {
		return nil, errors.New("slot does not belong to doctor")
	}

	if slot.Status() != availabilityDomain.SlotAvailable {
		return nil, errors.New("slot is not available")
	}

	if err := slot.MarkedBooked(now); err != nil {
		return nil, err
	}

	if err := uc.Slots.Update(ctx, tx, slot); err != nil {
		return nil, err
	}

	a, err := appointmentDomain.NewCreateAppointmentDomain(
		in.AppointmentID,
		in.DoctorID,
		in.PaymentID,
		in.SlotID,
		in.PriceCents,
		appointmentDomain.StatusScheduled,
		now,
	)

	if err != nil {
		return nil, err
	}
	if err := uc.Appointments.Create(ctx, tx, a); err != nil {
		return nil, err
	}

	provider, err := parseProvider(in.PaymentProvider)
	if err != nil {
		return nil, err
	}

	p, err := paymentDomain.NewPaymentDomain(
		in.PaymentID,
		in.AppointmentID,
		provider,
		in.PriceCents,
		paymentDomain.StatusPending,
		in.ExternalRef,
		now,
	)

	if err != nil {
		return nil, err
	}
	if err := uc.Payments.Create(ctx, tx, p); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.BookAppointmentOutput{
		AppointmentID: a.ID(),
		PaymentID:     p.ID(),
		Status:        string(a.Status()),
	}, nil

}

func parseProvider(raw string) (paymentDomain.Provider, error) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "STRIPE":
		return paymentDomain.ProviderStripe, nil
	case "MERCADOPAGO":
		return paymentDomain.ProviderMercadoPago, nil
	case "MANUAL":
		return paymentDomain.ProviderManual, nil
	default:
		return "", errors.New("invalid payment provider")
	}
}
