package config

import (
	"context"
	"fmt"
	"log"
	"time"
  

	// "github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	// Parse database config from env
	config, err := pgxpool.ParseConfig(DB_URL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	// Pool settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = time.Minute * 5

	// Initialize connection pool
	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	// Check database connection
	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL with pooling")
}
