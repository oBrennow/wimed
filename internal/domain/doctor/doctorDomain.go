package doctor

import (
	"errors"
	"wimed/internal/infra/http/restError"
)

type DoctorDomain struct {
	id     string
	name   string
	crm    string
	email  string
	active bool
}

func (d *DoctorDomain) Activate() error {
	if d.active {
		return errors.New("Doctor already activated")
	}
	d.active = true
	return nil
}

func (d *DoctorDomain) Deactivate() error {
	if !d.active {
		return errors.New("doctor is already inactive")
	}

	d.active = false
	return nil
}

func (d *DoctorDomain) SetID(id string) {
	d.id = id
}

func NewDomainDoctor(id string, name string, crm string, email string, status bool) (*DoctorDomain, error) {
	if name == "" {
		return nil, errors.New("doctor name is required")
	}
	if crm == "" {
		return nil, errors.New("doctor crm is required")
	}
	if email == "" {
		return nil, errors.New("doctor email is required")
	}

	return &DoctorDomain{
		id:     id,
		name:   name,
		crm:    crm,
		email:  email,
		active: status,
	}, nil
}

type DoctorDomainInterface interface {
	Activate() *restError.RestErr
	Deactivate() *restError.RestErr
	SetID(id string)
}
