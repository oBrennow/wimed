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

const (
	errPatientIDRequired     = "patient id is required"
	errPatientUserIDRequired = "user id is required"
	errPatientNameRequired   = "patient name is required"
)

func NewPatientDomain(id, userID, name string, active bool, now time.Time) (*PatientDomain, error) {
	id = strings.TrimSpace(id)
	userID = strings.TrimSpace(userID)
	name = strings.TrimSpace(name)

	if id == "" {
		return nil, errors.New(errPatientIDRequired)
	}
	if userID == "" {
		return nil, errors.New(errPatientUserIDRequired)
	}
	if name == "" {
		return nil, errors.New(errPatientNameRequired)
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

func RebuildPatientDomain(id, userID, name string, active bool, createdAt, updatedAt time.Time) *PatientDomain {
	return &PatientDomain{
		id:        id,
		userID:    userID,
		name:      name,
		active:    active,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (p *PatientDomain) ID() string           { return p.id }
func (p *PatientDomain) UserID() string       { return p.userID }
func (p *PatientDomain) Name() string         { return p.name }
func (p *PatientDomain) IsActive() bool       { return p.active }
func (p *PatientDomain) CreatedAt() time.Time { return p.createdAt }
func (p *PatientDomain) UpdatedAt() time.Time { return p.updatedAt }
