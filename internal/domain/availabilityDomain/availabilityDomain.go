package availabilityDomain

import (
	"errors"
	"strings"
	"time"
)

type SlotStatus string

var (
	ErrSlotIDRequired     = errors.New("slot: id is required")
	ErrSlotDoctorRequired = errors.New("slot: doctor_id is required")
	ErrSlotTimeInvalid    = errors.New("slot: time is invalid")
	ErrSlotStatusInvalid  = errors.New("slot: status is invalid")

	ErrSlotCannotBook           = errors.New("slot: can only book an available slot")
	ErrSlotCannotBlock          = errors.New("slot: can only block an available slot")
	ErrSlotCannotUnblock        = errors.New("slot: can only unblock a blocked slot")
	ErrSlotCreatedAtRequired    = errors.New("slot: created_at is required")
	ErrSlotUpdatedBeforeCreated = errors.New("slot: updated_at must be >= created_at")
)

const (
	SlotAvailable SlotStatus = "AVAILABLE"
	SlotBooked    SlotStatus = "BOOKED"
	SlotBlocked   SlotStatus = "BLOCKED"
)

func (st SlotStatus) IsValid() bool {
	switch st {
	case SlotAvailable, SlotBooked, SlotBlocked:
		return true
	default:
		return false
	}
}

type SlotDomain struct {
	id       string
	doctorID string

	startedAt time.Time
	endedAt   time.Time

	status SlotStatus

	createdAt time.Time
	updatedAt time.Time
}

func CreateSlotDomain(id, doctorID string, startedAt, endedAt time.Time, status SlotStatus, now time.Time) (*SlotDomain, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrSlotIDRequired
	}
	doctorID = strings.TrimSpace(doctorID)
	if doctorID == "" {
		return nil, ErrSlotDoctorRequired
	}
	if startedAt.IsZero() || endedAt.IsZero() || !startedAt.Before(endedAt) {
		return nil, ErrSlotTimeInvalid
	}
	if !status.IsValid() {
		return nil, ErrSlotStatusInvalid
	}
	if now.IsZero() {
		now = time.Now()
	}

	return &SlotDomain{
		id:        id,
		doctorID:  doctorID,
		startedAt: startedAt,
		endedAt:   endedAt,
		status:    status,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func RebuildSlotDomain(id, doctorID string, startedAt, endedAt time.Time, status SlotStatus, createdAt, updatedAt time.Time) (*SlotDomain, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrSlotIDRequired
	}
	doctorID = strings.TrimSpace(doctorID)
	if doctorID == "" {
		return nil, ErrSlotDoctorRequired
	}
	if startedAt.IsZero() || endedAt.IsZero() || !startedAt.Before(endedAt) {
		return nil, ErrSlotTimeInvalid
	}
	if !status.IsValid() {
		return nil, ErrSlotStatusInvalid
	}
	if createdAt.IsZero() {
		return nil, ErrSlotCreatedAtRequired
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	if updatedAt.Before(createdAt) {
		return nil, ErrSlotUpdatedBeforeCreated
	}

	return &SlotDomain{
		id:        id,
		doctorID:  doctorID,
		startedAt: startedAt,
		endedAt:   endedAt,
		status:    status,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (s *SlotDomain) ID() string           { return s.id }
func (s *SlotDomain) DoctorID() string     { return s.doctorID }
func (s *SlotDomain) StartedAt() time.Time { return s.startedAt }
func (s *SlotDomain) EndedAt() time.Time   { return s.endedAt }
func (s *SlotDomain) Status() SlotStatus   { return s.status }
func (s *SlotDomain) CreatedAt() time.Time { return s.createdAt }
func (s *SlotDomain) UpdatedAt() time.Time { return s.updatedAt }

func (s *SlotDomain) MarkBooked(now time.Time) error {
	if s.status != SlotAvailable {
		return ErrSlotCannotBook
	}
	s.status = SlotBooked
	s.touch(now)
	return nil
}

func (s *SlotDomain) Block(now time.Time) error {
	if s.status != SlotAvailable {
		return ErrSlotCannotBlock
	}
	s.status = SlotBlocked
	s.touch(now)
	return nil
}

func (s *SlotDomain) Unblock(now time.Time) error {
	if s.status != SlotBlocked {
		return ErrSlotCannotUnblock
	}
	s.status = SlotAvailable
	s.touch(now)
	return nil
}

func (s *SlotDomain) touch(now time.Time) {
	if now.IsZero() {
		now = time.Now()
	}
	s.updatedAt = now
}
