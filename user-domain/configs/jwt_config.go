package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadJWTSecret() (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secret := os.Getenv("USER_DOMAIN_JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT secret not found in environment variables")
	}

	return secret, nil
}