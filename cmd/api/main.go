package main

import (
	"fmt"
	"log"
	"os"
	"context"

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

	fmt.Println("Connected to database succesfully")
}
