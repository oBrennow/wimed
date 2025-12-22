package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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

	now := time.Now().UTC()

	// 1) user do médico e user do paciente
	_, err = pool.Exec(context.Background(), `
INSERT INTO users (id, email, password_hash, active, created_at, updated_at)
VALUES
  ('usr_doc_1','doc1@wimed.test','hash',true,$1,$1),
  ('usr_pat_1','pat1@wimed.test','hash',true,$1,$1)
ON CONFLICT (id) DO NOTHING;
`, now)
	if err != nil {
		log.Fatal(err)
	}

	// 2) doctor e patient
	_, err = pool.Exec(context.Background(), `
INSERT INTO doctors (id, user_id, name, registry_type, registry_number, specialty, session_minutes, price_cents, active, created_at, updated_at)
VALUES ('doc_1','usr_doc_1','Dr Test','CRM','12345','clinico',30,10000,true,$1,$1)
ON CONFLICT (id) DO NOTHING;
`, now)
	if err != nil {
		log.Fatal(err)
	}

	_, err = pool.Exec(context.Background(), `
INSERT INTO patients (id, user_id, name, active, created_at, updated_at)
VALUES ('pat_1','usr_pat_1','Paciente Teste',true,$1,$1)
ON CONFLICT (id) DO NOTHING;
`, now)
	if err != nil {
		log.Fatal(err)
	}

	// 3) slot disponível
	start := now.Add(24 * time.Hour)
	end := start.Add(30 * time.Minute)

	_, err = pool.Exec(context.Background(), `
INSERT INTO availability_slots (id, doctor_id, start_at, end_at, status, created_at, updated_at)
VALUES ('slot_1','doc_1',$1,$2,'AVAILABLE',$3,$3)
ON CONFLICT (id) DO NOTHING;
`, start, end, now)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Seed OK: usr_doc_1, usr_pat_1, doc_1, pat_1, slot_1")
}
