package config

import (
  "os"
  "log"
	"github.com/joho/godotenv"
)

var DB_URL   string
var BASE_URL string

func SetVariables() {
	// Always load a base .env file first (e.g., to read APP_ENV)
	_ = godotenv.Load(".env")

	env := os.Getenv("APP_ENV") // Can now be read from .env

	var envFile string
	if env == "production" {
		envFile = ".env.prod"
	} else {
		envFile = ".env.dev"
	}

	// Load environment-specific overrides
	if err := godotenv.Overload(envFile); err != nil {
		log.Fatalf("❌ Error loading env file [%s]: %v", envFile, err)
	} else {
		log.Printf("✅ Successfully loaded env file [%s]", envFile)
	}

	DB_URL = os.Getenv("DATABASE_URL")
	BASE_URL = os.Getenv("BASE_URL")
}
