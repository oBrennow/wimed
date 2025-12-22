package main

import (
	"context"
	"log"
	"os"

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

	sqlBytes, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		log.Fatalf("read migration: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("pgxpool.New: %v", err)
	}
	defer pool.Close()

	if _, err := pool.Exec(context.Background(), string(sqlBytes)); err != nil {
		log.Fatalf("apply migration: %v", err)
	}

	log.Println("Migration applied successfully.")
}
