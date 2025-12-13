package availabilityDomain

import (
	"errors"
	"time"
)

type SlotStatus string

const (
	SlotAvailable SlotStatus = "AVAILABLE"
	SlotBooked    SlotStatus = "BOOKED"
	SlotBlocked   SlotStatus = "BLOCKED"
)

type SlotDomain struct {
	id       string
	doctorID string

	startedAT time.Time
	endedAT   time.Time

	status SlotStatus

	createdAT time.Time
	updatedAT time.Time
}

const (
	errSlotIDRequired     = "slot id is required"
	errSlotDoctorRequired = "doctor id is required"
	errSlotTimeInvalid    = "slot start/end time is invalid"
	errSlotStatusInvalid  = "slot status is invalid"
)

func NewCreateSlotDomain(id, doctorID string, startAt, endAt time.Time, status SlotStatus, now time.Time) (*SlotDomain, error) {
	if id == "" {
		return nil, errors.New(errSlotIDRequired)
	}
	if doctorID == "" {
		return nil, errors.New(errSlotDoctorRequired)
	}
	if startAt.IsZero() || endAt.IsZero() || !startAt.Before(endAt) {
		return nil, errors.New(errSlotTimeInvalid)
	}
	if status != SlotAvailable && status != SlotBooked && status != SlotBlocked {
		return nil, errors.New(errSlotStatusInvalid)
	}
	if now.IsZero() {
		now = time.Now()
	}

	return &SlotDomain{
		id:        id,
		doctorID:  doctorID,
		startedAT: startAt,
		endedAT:   endAt,
		status:    status,
		createdAT: now,
		updatedAT: now,
	}, nil
}

func RebuildSlot(id, doctorID string, startAt, endAt time.Time, status SlotStatus, now time.Time) *SlotDomain {
	return &SlotDomain{
		id:        id,
		doctorID:  doctorID,
		startedAT: startAt,
		endedAT:   endAt,
		status:    status,
		createdAT: now,
		updatedAT: now,
	}
}

func (s *SlotDomain) ID() string           { return s.id }
func (s *SlotDomain) DoctorID() string     { return s.doctorID }
func (s *SlotDomain) StartedAt() time.Time { return s.startedAT }
func (s *SlotDomain) EndedAt() time.Time   { return s.endedAT }
func (s *SlotDomain) Status() SlotStatus   { return s.status }
func (s *SlotDomain) CreatedAt() time.Time { return s.createdAT }
func (s *SlotDomain) UpdatedAt() time.Time { return s.updatedAT }

func (s *SlotDomain) MarkedBooked(now time.Time) error {
	if s.status != SlotAvailable {
		return errors.New("slot is not available")
	}
	s.status = SlotBooked
	s.touch(now)
	return nil

}

func (s *SlotDomain) MarkAvailable(now time.Time) error {
	if s.status == SlotBooked {
		return errors.New("cannot revert booked slot to available directly")
	}
	s.status = SlotAvailable
	s.touch(now)
	return nil

}

func (s *SlotDomain) Block(now time.Time) error {
	if s.status == SlotBlocked {
		return errors.New("cannot block a booked slot")
	}
	s.status = SlotBlocked
	s.touch(now)
	return nil

}

func (s *SlotDomain) touch(now time.Time) {
	if now.IsZero() {
		now = time.Now()
	}
	s.updatedAT = now
}
