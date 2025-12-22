package appointmentDomain

import (
	"errors"
	"strings"
	"time"
)

type Status string

var (
	ErrApptIDRequired      = errors.New("appointment: id is required")
	ErrApptDoctorRequired  = errors.New("appointment: doctor is required")
	ErrApptPatientRequired = errors.New("appointment: patient id is required")
	ErrApptSlotRequired    = errors.New("appointment: slot id is required")
	ErrApptPriceInvalid    = errors.New("appointment: price cents must be >= 0")
	ErrApptStatusInvalid   = errors.New("appointment: status is invalid")

	ErrOnlyPaidCanComplete      = errors.New("appointment: only paid can complete")
	ErrApptCreatedAtRequired    = errors.New("appointment: created_at is required")
	ErrApptUpdatedBeforeCreated = errors.New("appointment: updated_at must be >= created_at")
	ErrApptNotScheduled         = errors.New("appointment: must be scheduled")
	ErrApptCannotCancel         = errors.New("appointment: cannot cancel in current status")
)

const (
	StatusScheduled Status = "SCHEDULED"
	StatusPaid      Status = "PAID"
	StatusCanceled  Status = "CANCELED"
	StatusCompleted Status = "COMPLETED"
)

func (st Status) IsValid() bool {
	switch st {
	case StatusScheduled, StatusPaid, StatusCanceled, StatusCompleted:
		return true
	default:
		return false
	}
}

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

func CreateAppointmentDomain(id, patientID, doctorID, slotID string, priceCents int64, status Status, now time.Time) (*AppointmentDomain, error) {
	id = strings.TrimSpace(id)
	patientID = strings.TrimSpace(patientID)
	doctorID = strings.TrimSpace(doctorID)
	slotID = strings.TrimSpace(slotID)

	if id == "" {
		return nil, ErrApptIDRequired
	}
	if doctorID == "" {
		return nil, ErrApptDoctorRequired
	}
	if patientID == "" {
		return nil, ErrApptPatientRequired
	}
	if slotID == "" {
		return nil, ErrApptSlotRequired
	}
	if priceCents < 0 {
		return nil, ErrApptPriceInvalid
	}
	if !status.IsValid() {
		return nil, ErrApptStatusInvalid
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

func RebuildAppointmentDomain(id, patientID, doctorID, slotID string, priceCents int64, status Status, createdAt, updatedAt time.Time) (*AppointmentDomain, error) {
	id = strings.TrimSpace(id)
	patientID = strings.TrimSpace(patientID)
	doctorID = strings.TrimSpace(doctorID)
	slotID = strings.TrimSpace(slotID)

	if id == "" {
		return nil, ErrApptIDRequired
	}
	if doctorID == "" {
		return nil, ErrApptDoctorRequired
	}
	if patientID == "" {
		return nil, ErrApptPatientRequired
	}
	if slotID == "" {
		return nil, ErrApptSlotRequired
	}
	if priceCents < 0 {
		return nil, ErrApptPriceInvalid
	}
	if !status.IsValid() {
		return nil, ErrApptStatusInvalid
	}
	if createdAt.IsZero() {
		return nil, ErrApptCreatedAtRequired
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	if updatedAt.Before(createdAt) {
		return nil, ErrApptUpdatedBeforeCreated
	}
	return &AppointmentDomain{
		id:         id,
		doctorID:   doctorID,
		patientID:  patientID,
		slotID:     slotID,
		priceCents: priceCents,
		status:     status,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}, nil
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
		return ErrApptNotScheduled
	}
	a.status = StatusPaid
	a.touch(now)
	return nil
}

func (a *AppointmentDomain) Cancel(now time.Time) error {
	if a.status == StatusCanceled || a.status == StatusCompleted {
		return ErrApptCannotCancel
	}
	a.status = StatusCanceled
	a.touch(now)
	return nil
}

func (a *AppointmentDomain) Complete(now time.Time) error {
	if a.status != StatusPaid {
		return ErrOnlyPaidCanComplete
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
