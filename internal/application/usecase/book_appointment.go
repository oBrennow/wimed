package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"wimed/internal/application/dto"
	"wimed/internal/domain/appointmentDomain"
	"wimed/internal/domain/paymentDomain"

	"wimed/internal/application/ports"
)

var (
	ErrBookAppointmentAppointmentIDRequired = errors.New("book_appointment: appointment_id is required")
	ErrBookAppointmentPatientIDRequired     = errors.New("book_appointment: patient_id is required")
	ErrBookAppointmentSlotIDRequired        = errors.New("book_appointment: slot_id is required")
	ErrBookAppointmentDoctorIDRequired      = errors.New("book_appointment: doctor_id is required")
	ErrBookAppointmentPaymentIDRequired     = errors.New("book_appointment: payment_id is required")
	ErrBookAppointmentPriceInvalid          = errors.New("book_appointment: price_cents must be >= 0")

	ErrBookAppointmentPatientNotFound        = errors.New("book_appointment: patient not found")
	ErrBookAppointmentSlotDoctorMismatch     = errors.New("book_appointment: slot does not belong to doctor")
	ErrBookAppointmentInvalidPaymentProvider = errors.New("book_appointment: invalid payment provider")
)

type BookAppointment struct {
	TxManager ports.TxManager

	Patients     ports.PatientRepository
	Slots        ports.SlotLockRepository
	Appointments ports.AppointmentRepository
	Payments     ports.PaymentRepository

	Now func() time.Time
}

func (uc *BookAppointment) Execute(ctx context.Context, in dto.BookAppointmentInput) (*dto.BookAppointmentOutput, error) {
	appointmentID := strings.TrimSpace(in.AppointmentID)
	patientID := strings.TrimSpace(in.PatientID)
	slotID := strings.TrimSpace(in.SlotID)
	doctorID := strings.TrimSpace(in.DoctorID)
	paymentID := strings.TrimSpace(in.PaymentID)

	if appointmentID == "" {
		return nil, ErrBookAppointmentAppointmentIDRequired
	}
	if patientID == "" {
		return nil, ErrBookAppointmentPatientIDRequired
	}
	if slotID == "" {
		return nil, ErrBookAppointmentSlotIDRequired
	}
	if doctorID == "" {
		return nil, ErrBookAppointmentDoctorIDRequired
	}
	if paymentID == "" {
		return nil, ErrBookAppointmentPaymentIDRequired
	}
	if in.PriceCents < 0 {
		return nil, ErrBookAppointmentPriceInvalid
	}

	now := time.Now()
	if uc.Now != nil {
		now = uc.Now()
	}

	provider, err := parseProvider(in.PaymentProvider)
	if err != nil {

		return nil, fmt.Errorf("%w: %v", ErrBookAppointmentInvalidPaymentProvider, err)
	}

	tx, err := uc.TxManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	ok, err := uc.Patients.ExistsByID(ctx, tx, patientID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrBookAppointmentPatientNotFound
	}

	slot, err := uc.Slots.GetByIDForUpdate(ctx, tx, slotID)
	if err != nil {
		return nil, err
	}

	if slot.DoctorID() != doctorID {
		return nil, ErrBookAppointmentSlotDoctorMismatch
	}

	if err := slot.MarkBooked(now); err != nil {
		// domínio já devolve ErrSlotCannotBook etc.
		return nil, err
	}

	if err := uc.Slots.Update(ctx, tx, slot); err != nil {
		return nil, err
	}

	appt, err := appointmentDomain.CreateAppointmentDomain(
		appointmentID,
		patientID,
		slot.DoctorID(),
		slotID,
		in.PriceCents,
		appointmentDomain.StatusScheduled,
		now,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.Appointments.Create(ctx, tx, appt); err != nil {
		return nil, err
	}

	pay, err := paymentDomain.CreatePaymentDomain(
		paymentID,
		appointmentID,
		provider,
		in.PriceCents,
		paymentDomain.StatusPending,
		in.ExternalRef,
		now,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.Payments.Create(ctx, tx, pay); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.BookAppointmentOutput{
		AppointmentID: appt.ID(),
		PaymentID:     pay.ID(),
		Status:        string(appt.Status()),
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
		return "", fmt.Errorf("provider=%q", raw)
	}
}
