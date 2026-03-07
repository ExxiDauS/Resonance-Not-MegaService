package database

import (
	"log"

	"match-domain/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDatabaseClient() (*gorm.DB, error) {
	log.Printf("===============================")
	log.Printf("Loading database configuration")
	dbConfig, err := configs.LoadDatabaseConfig()
	log.Printf("Database configuration loaded")
	if err != nil {
		log.Fatalf("Failed to load database configuration: %v", err)
		return nil, err
	}
	log.Printf("===============================")
	log.Printf("Connecting to Postgres database")
	dbURL := dbConfig.DbURL
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to Postgres database: %v", err)
		return nil, err
	}
	log.Printf("Connected to Postgres database")
	log.Printf("===============================")
	return db, nil
}
