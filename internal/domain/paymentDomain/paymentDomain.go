package paymentDomain

import (
	"errors"
	"time"
)

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusApproved  Status = "APPROVED"
	StatusRejected  Status = "REJECTED"
	StatusRefunded  Status = "REFUNDED"
	StatusCancelled Status = "CANCELLED"
)

type Provider string

const (
	ProviderStripe      Provider = "STRIPE"
	ProviderMercadoPago Provider = "MERCADOPAGO"
	ProviderManual      Provider = "MANUAL"
)

type PaymentDomain struct {
	id            string
	appointmentID string

	provider    Provider
	amountCents int64

	status      Status
	externalRef string

	createdAt time.Time
	updatedAt time.Time
}

const (
	errPayIDRequired      = "PaymentDomain id is required"
	errPayApptRequired    = "appointment id is required"
	errPayAmountInvalid   = "amount cents must be >= 0"
	errPayProviderInvalid = "PaymentDomain provider is invalid"
	errPayStatusInvalid   = "PaymentDomain status is invalid"
)

func NewPaymentDomain(
	id, appointmentID string,
	provider Provider,
	amountCents int64,
	status Status,
	externalRef string,
	now time.Time,
) (*PaymentDomain, error) {
	if id == "" {
		return nil, errors.New(errPayIDRequired)
	}
	if appointmentID == "" {
		return nil, errors.New(errPayApptRequired)
	}
	if amountCents < 0 {
		return nil, errors.New(errPayAmountInvalid)
	}
	if provider != ProviderStripe && provider != ProviderMercadoPago && provider != ProviderManual {
		return nil, errors.New(errPayProviderInvalid)
	}
	if status != StatusPending && status != StatusApproved && status != StatusRejected && status != StatusRefunded && status != StatusCancelled {
		return nil, errors.New(errPayStatusInvalid)
	}
	if now.IsZero() {
		now = time.Now()
	}
	return &PaymentDomain{
		id:            id,
		appointmentID: appointmentID,
		provider:      provider,
		amountCents:   amountCents,
		status:        status,
		externalRef:   externalRef,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

func (p *PaymentDomain) ID() string            { return p.id }
func (p *PaymentDomain) AppointmentID() string { return p.appointmentID }
func (p *PaymentDomain) Provider() Provider    { return p.provider }
func (p *PaymentDomain) AmountCents() int64    { return p.amountCents }
func (p *PaymentDomain) Status() Status        { return p.status }
func (p *PaymentDomain) ExternalRef() string   { return p.externalRef }
func (p *PaymentDomain) CreatedAt() time.Time  { return p.createdAt }
func (p *PaymentDomain) UpdatedAt() time.Time  { return p.updatedAt }

func (p *PaymentDomain) Approve(now time.Time) error {
	if p.status != StatusApproved {
		return errors.New("payment is not pending")
	}
	p.status = StatusApproved
	p.touch(now)
	return nil
}

func (p *PaymentDomain) Reject(now time.Time) error {
	if p.status != StatusPending {
		return errors.New("payment is not pending")
	}
	p.status = StatusRejected
	p.touch(now)
	return nil

}

func (p *PaymentDomain) touch(now time.Time) {
	if now.IsZero() {
		now = time.Now()
	}
}
