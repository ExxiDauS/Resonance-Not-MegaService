package configs

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DbURL string
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("INTERACTION_DOMAIN_URL")
	username := os.Getenv("INTERACTION_DOMAIN_USERNAME")
	password := os.Getenv("INTERACTION_DOMAIN_PASSWORD")

	// Remove JDBC prefix if present
	dbURL = strings.TrimPrefix(dbURL, "jdbc:")
	// Remove postgresql:// prefix to rebuild the connection string
	dbURL = strings.TrimPrefix(dbURL, "postgresql://")

	// Construct proper PostgreSQL connection string with credentials
	dbURL = fmt.Sprintf("postgresql://%s:%s@%s?sslmode=disable", username, password, dbURL)

	config := &DatabaseConfig{
		DbURL: dbURL,
	}

	return config, nil
}
