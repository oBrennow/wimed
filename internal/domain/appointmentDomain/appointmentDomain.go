package appointmentDomain

import (
	"errors"
	"time"
)

type Status string

const (
	StatusScheduled Status = "SCHEDULED"
	StatusPaid      Status = "PAID"
	StatusCanceled  Status = "CANCELED"
	StatusCompleted Status = "COMPLETED"
)

type AppointmentDomain struct {
	id        string
	doctorID  string
	patientID string
	slotID    string

	priceCents int64

	status Status

	createdAt time.Time
	updatedAt time.Time
}

const (
	errApptIDRequired      = "appointment id is required"
	errApptDoctorRequired  = "doctor id is required"
	errApptPatientRequired = "patient id is required"
	errApptSlotRequired    = "slot id is required"
	errApptPriceInvalid    = "price cents must be >= 0"
	errApptStatusInvalid   = "appointment status is invalid"
)

func NewCreateAppointmentDomain(id, patientID, doctorID, slotID string, priceCents int64, status Status, now time.Time) (*AppointmentDomain, error) {
	if id == "" {
		return nil, errors.New(errApptIDRequired)
	}
	if doctorID == "" {
		return nil, errors.New(errApptDoctorRequired)
	}
	if patientID == "" {
		return nil, errors.New(errApptPatientRequired)
	}
	if slotID == "" {
		return nil, errors.New(errApptSlotRequired)
	}
	if priceCents < 0 {
		return nil, errors.New(errApptPriceInvalid)
	}
	if status != StatusScheduled && status != StatusPaid && status != StatusCanceled && status != StatusCompleted {
		return nil, errors.New(errApptStatusInvalid)
	}

	if now.IsZero() {
		now = time.Now()
	}

	return &AppointmentDomain{
		id:         id,
		doctorID:   doctorID,
		patientID:  patientID,
		slotID:     slotID,
		priceCents: priceCents,
		status:     status,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

func RebuildAppointment(id, patientID, doctorID, slotID string, priceCents int64, status Status, now time.Time) *AppointmentDomain {
	return &AppointmentDomain{
		id:         id,
		doctorID:   doctorID,
		patientID:  patientID,
		slotID:     slotID,
		priceCents: priceCents,
		status:     status,
		createdAt:  now,
		updatedAt:  now,
	}
}

func (a *AppointmentDomain) ID() string           { return a.id }
func (a *AppointmentDomain) DoctorID() string     { return a.doctorID }
func (a *AppointmentDomain) PatientID() string    { return a.patientID }
func (a *AppointmentDomain) SlotID() string       { return a.slotID }
func (a *AppointmentDomain) PriceCents() int64    { return a.priceCents }
func (a *AppointmentDomain) Status() Status       { return a.status }
func (a *AppointmentDomain) CreatedAt() time.Time { return a.createdAt }
func (a *AppointmentDomain) UpdatedAt() time.Time { return a.updatedAt }

func (a *AppointmentDomain) MarkPaid(now time.Time) error {
	if a.status != StatusScheduled {
		return errors.New("appointment is not scheduled status")
	}
	a.status = StatusPaid
	a.touch(now)
	return nil
}

func (a *AppointmentDomain) Cancel(now time.Time) error {
	if a.status == StatusCanceled || a.status == StatusCompleted {
		return errors.New("appointment cannot be canceled in current status")
	}
	a.status = StatusCanceled
	a.touch(now)
	return nil
}

func (a *AppointmentDomain) Complete(now time.Time) error {
	if a.status != StatusScheduled {
		return errors.New("only paid appointment can be completed")
	}
	a.status = StatusCompleted
	a.touch(now)
	return nil
}

func (a *AppointmentDomain) touch(now time.Time) {
	if now.IsZero() {
		now = time.Now()
	}
	a.updatedAt = now
}
