package patientDomain

import (
	"errors"
	"strings"
	"time"
)

type PatientDomain struct {
	id     string
	userID string
	name   string

	active    bool
	createdAt time.Time
	updatedAt time.Time
}

var (
	ErrPatientIDRequired           = errors.New("patient: id is required")
	ErrPatientUserIDRequired       = errors.New("patient: user_id is required")
	ErrPatientNameRequired         = errors.New("patient: name is required")
	ErrPatientAlreadyActive        = errors.New("patient: already active")
	ErrPatientAlreadyInactive      = errors.New("patient: already inactive")
	ErrPatientCreatedAtRequired    = errors.New("patient: created_at is required")
	ErrPatientUpdatedBeforeCreated = errors.New("patient: updated_at must be >= created_at")
)

func CreatePatientDomain(id, userID, name string, active bool, now time.Time) (*PatientDomain, error) {
	id = strings.TrimSpace(id)
	userID = strings.TrimSpace(userID)
	name = strings.TrimSpace(name)

	if id == "" {
		return nil, ErrPatientIDRequired
	}
	if userID == "" {
		return nil, ErrPatientUserIDRequired
	}
	if name == "" {
		return nil, ErrPatientNameRequired
	}

	if now.IsZero() {
		now = time.Now()
	}

	return &PatientDomain{
		id:        id,
		userID:    userID,
		name:      name,
		active:    active,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func RebuildPatientDomain(id, userID, name string, active bool, createdAt, updatedAt time.Time) (*PatientDomain, error) {
	id = strings.TrimSpace(id)
	userID = strings.TrimSpace(userID)
	name = strings.TrimSpace(name)

	if id == "" {
		return nil, ErrPatientIDRequired
	}
	if userID == "" {
		return nil, ErrPatientUserIDRequired
	}
	if name == "" {
		return nil, ErrPatientNameRequired
	}

	if createdAt.IsZero() {
		return nil, ErrPatientCreatedAtRequired
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	if updatedAt.Before(createdAt) {
		return nil, ErrPatientUpdatedBeforeCreated
	}

	return &PatientDomain{
		id:        id,
		userID:    userID,
		name:      name,
		active:    active,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (p *PatientDomain) ID() string           { return p.id }
func (p *PatientDomain) UserID() string       { return p.userID }
func (p *PatientDomain) Name() string         { return p.name }
func (p *PatientDomain) IsActive() bool       { return p.active }
func (p *PatientDomain) CreatedAt() time.Time { return p.createdAt }
func (p *PatientDomain) UpdatedAt() time.Time { return p.updatedAt }

func (p *PatientDomain) Activate(now time.Time) error {
	if p.active {
		return ErrPatientAlreadyActive
	}
	p.active = true
	p.touch(now)
	return nil
}

func (p *PatientDomain) Deactivate(now time.Time) error {
	if !p.active {
		return ErrPatientAlreadyInactive
	}
	p.active = false
	p.touch(now)
	return nil
}

func (p *PatientDomain) ChangeName(name string, now time.Time) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrPatientNameRequired
	}
	p.name = name
	p.touch(now)
	return nil
}

func (p *PatientDomain) touch(now time.Time) {
	if now.IsZero() {
		now = time.Now()
	}
	p.updatedAt = now
}
