package doctorDomain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrDoctorIDRequired           = errors.New("doctor: id is required")
	ErrDoctorUserIDRequired       = errors.New("doctor: user_id is required")
	ErrDoctorNameRequired         = errors.New("doctor: name is required")
	ErrDoctorRegistryInvalid      = errors.New("doctor: registry is invalid")
	ErrDoctorSessionInvalid       = errors.New("doctor: session_minutes must be > 0")
	ErrDoctorPriceInvalid         = errors.New("doctor: price_cents must be >= 0")
	ErrDoctorAlreadyActive        = errors.New("doctor: already active")
	ErrDoctorAlreadyInactive      = errors.New("doctor: already inactive")
	ErrDoctorCreatedAtRequired    = errors.New("doctor: created_at is required")
	ErrDoctorUpdatedBeforeCreated = errors.New("doctor: updated_at must be >= created_at")
)

type RegistryType string

const (
	RegistryCRM   RegistryType = "CRM"
	RegistryCRP   RegistryType = "CRP"
	RegistryOther RegistryType = "OTHER"
)

type DoctorDomain struct {
	id             string
	userID         string
	name           string
	registryType   RegistryType
	registryNumber string
	specialty      string
	sessionMinutes int
	priceCents     int64
	active         bool
	createdAt      time.Time
	updatedAt      time.Time
}

func CreateDoctorDomain(
	id string,
	userID string,
	name string,
	regType RegistryType,
	regNumber string,
	specialty string,
	sessionMinutes int,
	priceCents int64,
	active bool,
	now time.Time,
) (*DoctorDomain, error) {

	id = strings.TrimSpace(id)
	userID = strings.TrimSpace(userID)
	name = strings.TrimSpace(name)
	regNumber = strings.TrimSpace(regNumber)
	specialty = strings.TrimSpace(specialty)

	if id == "" {
		return nil, ErrDoctorIDRequired
	}
	if userID == "" {
		return nil, ErrDoctorUserIDRequired
	}
	if name == "" {
		return nil, ErrDoctorNameRequired
	}
	if !regType.IsValid() || regNumber == "" {
		return nil, ErrDoctorRegistryInvalid
	}
	if sessionMinutes <= 0 {
		return nil, ErrDoctorSessionInvalid
	}
	if priceCents < 0 {
		return nil, ErrDoctorPriceInvalid
	}
	if now.IsZero() {
		now = time.Now()
	}
	return &DoctorDomain{
		id:             id,
		userID:         userID,
		name:           name,
		registryType:   regType,
		registryNumber: regNumber,
		specialty:      specialty,
		sessionMinutes: sessionMinutes,
		priceCents:     priceCents,
		active:         active,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

func RebuildDoctorDomain(
	id, userID, name string,
	regType RegistryType, regNumber, specialty string,
	sessionMinutes int, priceCents int64,
	active bool,
	createdAt, updatedAt time.Time,
) (*DoctorDomain, error) {
	id = strings.TrimSpace(id)
	userID = strings.TrimSpace(userID)
	name = strings.TrimSpace(name)
	regNumber = strings.TrimSpace(regNumber)
	specialty = strings.TrimSpace(specialty)

	if id == "" {
		return nil, ErrDoctorIDRequired
	}
	if userID == "" {
		return nil, ErrDoctorUserIDRequired
	}
	if name == "" {
		return nil, ErrDoctorNameRequired
	}
	if !regType.IsValid() || regNumber == "" {
		return nil, ErrDoctorRegistryInvalid
	}
	if sessionMinutes <= 0 {
		return nil, ErrDoctorSessionInvalid
	}
	if priceCents < 0 {
		return nil, ErrDoctorPriceInvalid
	}
	if createdAt.IsZero() {
		return nil, ErrDoctorCreatedAtRequired
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	if updatedAt.Before(createdAt) {
		return nil, ErrDoctorUpdatedBeforeCreated
	}

	return &DoctorDomain{
		id:             id,
		userID:         userID,
		name:           name,
		registryType:   regType,
		registryNumber: regNumber,
		specialty:      specialty,
		sessionMinutes: sessionMinutes,
		priceCents:     priceCents,
		active:         active,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}, nil
}

func (d *DoctorDomain) ID() string           { return d.id }
func (d *DoctorDomain) UserID() string       { return d.userID }
func (d *DoctorDomain) Name() string         { return d.name }
func (d *DoctorDomain) Specialty() string    { return d.specialty }
func (d *DoctorDomain) SessionMinutes() int  { return d.sessionMinutes }
func (d *DoctorDomain) PriceCents() int64    { return d.priceCents }
func (d *DoctorDomain) IsActive() bool       { return d.active }
func (d *DoctorDomain) CreatedAt() time.Time { return d.createdAt }
func (d *DoctorDomain) UpdatedAt() time.Time { return d.updatedAt }

func (d *DoctorDomain) Activate(now time.Time) error {
	if d.active {
		return ErrDoctorAlreadyActive
	}
	d.active = true
	d.touch(now)
	return nil
}

func (d *DoctorDomain) Deactivate(now time.Time) error {
	if !d.active {
		return ErrDoctorAlreadyInactive
	}
	d.active = false
	d.touch(now)
	return nil
}

func (d *DoctorDomain) UpdatePricing(sessionMinutes int, priceCents int64, now time.Time) error {
	if sessionMinutes <= 0 {
		return ErrDoctorSessionInvalid
	}
	if priceCents < 0 {
		return ErrDoctorPriceInvalid
	}
	d.sessionMinutes = sessionMinutes
	d.priceCents = priceCents
	d.touch(now)
	return nil
}
func (t RegistryType) IsValid() bool {
	switch t {
	case RegistryCRM, RegistryCRP, RegistryOther:
		return true
	default:
		return false
	}
}

func (d *DoctorDomain) touch(now time.Time) {
	if now.IsZero() {
		now = time.Now()
	}
	d.updatedAt = now
}
