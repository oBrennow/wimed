package main

import (
	"context"


	"log"
	"os"
	"time"

	apphttp "wimed/internal/infra/http"
	"wimed/internal/application/usecase"
	"wimed/internal/infra/http/handlers"
	"wimed/internal/infra/persistence/postgres"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)


func main() { 
	env := os.Getenv("ENV")
	if env == ""{
		env = "dev"
	}

	if env != "production" {
		if err := godotenv.Load(); err != nil {
			log.Printf("warning: .env not loaded: %v", err)
		}
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("pgxpool.New error: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("db ping error %v", err)
	}

	bookUC := &usecase.BookAppointment{
		TxManager: postgres.TxManager{Pool: pool},
		Patients: postgres.PatientRepository{},
		Slots: postgres.SlotRepository{},
		Appointments: postgres.AppointmentRepository{},
		Payments: postgres.PaymentRepository{},
		Now: func()time.Time {return time.Now().UTC()},
	}

	appointmentHandler := handlers.NewAppointmentHandler(bookUC)

	listSlotsUC := &usecase.ListAvailableSlots{
		TxManager: postgres.TxManager{Pool: pool},
		Slots:     postgres.SlotRepository{},
		Now:       func() time.Time { return time.Now().UTC() },
	}

	slotHandler := handlers.NewSlotHandler(listSlotsUC)


	r := apphttp.NewRouter(appointmentHandler, slotHandler)

	gin.SetMode(gin.ReleaseMode)

	log.Println("listening on: 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
