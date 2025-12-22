package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"wimed/internal/application/dto"
	"wimed/internal/application/usecase"
	"wimed/internal/infra/persistence/postgres"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("could not load .env: %v", err)
	} else {
		log.Printf(".env loaded")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("pgxpool.New: %v", err)
	}
	defer pool.Close()

	uc := &usecase.BookAppointment{
		TxManager: postgres.TxManager{Pool: pool},

		Patients:     postgres.PatientRepository{},
		Slots:        postgres.SlotRepository{},
		Appointments: postgres.AppointmentRepository{},
		Payments:     postgres.PaymentRepository{},

		Now: func() time.Time { return time.Now().UTC() },
	}

	out, err := uc.Execute(context.Background(), dto.BookAppointmentInput{
		AppointmentID: "appt_1",
		PaymentID:     "pay_1",
		SlotID:        "slot_1",
		DoctorID:      "doc_1",
		PatientID:     "pat_1",

		PriceCents:      10000,
		PaymentProvider: "MERCADOPAGO",
		ExternalRef:     "ext_test_1",
	})
	if err != nil {
		log.Fatalf("BookAppointment error: %v", err)
	}

	log.Printf("OK: appointment=%s payment=%s status=%s\n", out.AppointmentID, out.PaymentID, out.Status)
}
