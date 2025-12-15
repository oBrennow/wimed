package dto

type BookAppointmentInput struct {
	AppointmentID	string
	PaymentID		string

	SlotID			string
	DoctorID		string
	PatientID		string

	PriceCents		string
	PaymentProvider	string
	ExternalRef		string
}

type BookAppointmentOutput struct {
	Appointment	string
	PaymentID	string
	Status		string
}