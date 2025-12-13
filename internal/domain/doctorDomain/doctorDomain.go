package doctorDomain

import (
	"errors"
	"strings"
	"time"
)

type RegistryType string

const (
	RegistryCRM   RegistryType = "CRM"
	RegistryCRP   RegistryType = "CRP"
	RegistryOther RegistryType = "other"
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

const (
	errDoctorIDRequired      = "doctor id is required"
	errDoctorUserIDRequired  = "user id is required"
	errDoctorNameRequired    = "doctor name is required"
	errDoctorRegistryInvalid = "registry type/number is invalid"
	errDoctorSessionInvalid  = "session minutes must be > 0"
	errDoctorPriceInvalid    = "price cents must be >= 0"
)

func NewCreateDoctorDomain(
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
		return nil, errors.New(errDoctorIDRequired)
	}
	if userID == "" {
		return nil, errors.New(errDoctorUserIDRequired)
	}
	if name == "" {
		return nil, errors.New(errDoctorNameRequired)
	}
	if (regType != RegistryCRM && regType != RegistryCRP && regType != RegistryOther) || regNumber == "" {
		return nil, errors.New(errDoctorRegistryInvalid)
	}
	if sessionMinutes <= 0 {
		return nil, errors.New(errDoctorSessionInvalid)
	}
	if priceCents < 0 {
		return nil, errors.New(errDoctorPriceInvalid)
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
	regType RegistryType, regNumber, speciality string,
	sessionMinutes int, priceCents int64,
	active bool,
	createdAt, updatedAt time.Time,
) *DoctorDomain {
	return &DoctorDomain{
		id:             id,
		userID:         userID,
		name:           name,
		registryType:   regType,
		registryNumber: regNumber,
		specialty:      speciality,
		sessionMinutes: sessionMinutes,
		priceCents:     priceCents,
		active:         active,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

func (d *DoctorDomain) ID() string           { return d.id }
func (d *DoctorDomain) UserID() string       { return d.userID }
func (d *DoctorDomain) Name() string         { return d.name }
func (d *DoctorDomain) Speciality() string   { return d.specialty }
func (d *DoctorDomain) SessionMinutes() int  { return d.sessionMinutes }
func (d *DoctorDomain) PriceCents() int64    { return d.priceCents }
func (d *DoctorDomain) IsActive() bool       { return d.active }
func (d *DoctorDomain) CreatedAt() time.Time { return d.createdAt }
func (d *DoctorDomain) UpdatedAt() time.Time { return d.updatedAt }

func (d *DoctorDomain) Activate(now time.Time) error {
	if d.active {
		return errors.New("userDomain is already active")
	}
	d.active = true
	d.touch(now)
	return nil
}

func (d *DoctorDomain) Deactivate(now time.Time) error {
	if !d.active {
		return errors.New("userDomain is already inactive")
	}
	d.active = false
	d.touch(now)
	return nil
}

func (d *DoctorDomain) UpdatePricing(sessionMinutes int, priceCents int64, now time.Time) error {
	if sessionMinutes <= 0 {
		return errors.New(errDoctorPriceInvalid)
	}
	if priceCents < 0 {
		return errors.New(errDoctorPriceInvalid)
	}
	d.sessionMinutes = sessionMinutes
	d.priceCents = priceCents
	d.touch(now)
	return nil
}

func (d *DoctorDomain) touch(now time.Time) {
	if now.IsZero() {
		now = time.Now()
	}
	d.updatedAt = now
}
