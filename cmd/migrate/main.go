package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
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
