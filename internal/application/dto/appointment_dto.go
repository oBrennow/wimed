package dto

type BookAppointmentInput struct {
	AppointmentID string
	PaymentID     string

	SlotID    string
	DoctorID  string
	PatientID string

	PriceCents      int64
	PaymentProvider string
	ExternalRef     string
}

type BookAppointmentOutput struct {
	AppointmentID string
	PaymentID     string
	Status        string
}
