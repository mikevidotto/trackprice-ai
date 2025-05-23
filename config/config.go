package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadConfig reads .env file
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
