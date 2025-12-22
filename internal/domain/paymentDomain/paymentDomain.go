package paymentDomain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrPayIDRequired      = errors.New("payment: id is required")
	ErrPayApptRequired    = errors.New("payment: appointment_id is required")
	ErrPayAmountInvalid   = errors.New("payment: amount_cents must be >= 0")
	ErrPayProviderInvalid = errors.New("payment: provider is invalid")
	ErrPayStatusInvalid   = errors.New("payment: status is invalid")

	ErrPaymentIsNotPending     = errors.New("payment: payment is not pending")
	ErrPaymentIsNotApproved    = errors.New("payment: payment is not approved")
	ErrPayCreatedAtRequired    = errors.New("payment: created_at is required")
	ErrPayUpdatedBeforeCreated = errors.New("payment: updated_at must be >= created_at")
)

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusApproved  Status = "APPROVED"
	StatusRejected  Status = "REJECTED"
	StatusRefunded  Status = "REFUNDED"
	StatusCancelled Status = "CANCELLED"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusApproved, StatusRefunded, StatusCancelled, StatusRejected:
		return true
	default:
		return false
	}
}

type Provider string

const (
	ProviderStripe      Provider = "STRIPE"
	ProviderMercadoPago Provider = "MERCADOPAGO"
	ProviderManual      Provider = "MANUAL"
)

func (p Provider) IsValid() bool {
	switch p {
	case ProviderStripe, ProviderMercadoPago, ProviderManual:
		return true
	default:
		return false
	}
}

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

func CreatePaymentDomain(id, appointmentID string, provider Provider, amountCents int64, status Status, externalRef string, now time.Time) (*PaymentDomain, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrPayIDRequired
	}
	appointmentID = strings.TrimSpace(appointmentID)
	if appointmentID == "" {
		return nil, ErrPayApptRequired
	}
	if amountCents < 0 {
		return nil, ErrPayAmountInvalid
	}
	if !provider.IsValid() {
		return nil, ErrPayProviderInvalid
	}
	if !status.IsValid() {
		return nil, ErrPayStatusInvalid
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

func RebuildPaymentDomain(id, appointmentID string, provider Provider, amountCents int64, status Status, externalRef string, createdAt, updatedAt time.Time) (*PaymentDomain, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrPayIDRequired
	}
	appointmentID = strings.TrimSpace(appointmentID)
	if appointmentID == "" {
		return nil, ErrPayApptRequired
	}
	if amountCents < 0 {
		return nil, ErrPayAmountInvalid
	}
	if !provider.IsValid() {
		return nil, ErrPayProviderInvalid
	}
	if !status.IsValid() {
		return nil, ErrPayStatusInvalid
	}
	if createdAt.IsZero() {
		return nil, ErrPayCreatedAtRequired
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	if updatedAt.Before(createdAt) {
		return nil, ErrPayUpdatedBeforeCreated
	}

	return &PaymentDomain{
		id:            id,
		appointmentID: appointmentID,
		provider:      provider,
		amountCents:   amountCents,
		status:        status,
		externalRef:   externalRef,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
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
	if p.status != StatusPending {
		return ErrPaymentIsNotPending
	}
	p.status = StatusApproved
	p.touch(now)
	return nil
}

func (p *PaymentDomain) Reject(now time.Time) error {
	if p.status != StatusPending {
		return ErrPaymentIsNotPending
	}
	p.status = StatusRejected
	p.touch(now)
	return nil

}

func (p *PaymentDomain) Cancel(now time.Time) error {
	if p.status != StatusPending {
		return ErrPaymentIsNotPending
	}
	p.status = StatusCancelled
	p.touch(now)
	return nil
}

func (p *PaymentDomain) Refund(now time.Time) error {
	if p.status != StatusApproved {
		return ErrPaymentIsNotApproved
	}
	p.status = StatusRefunded
	p.touch(now)
	return nil
}

func (p *PaymentDomain) touch(now time.Time) {
	if now.IsZero() {
		now = time.Now()
	}
	p.updatedAt = now
}
