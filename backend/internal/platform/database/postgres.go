package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB   *gorm.DB
	once sync.Once
)

// InitDB initializes the database connection
func InitDB() {
	once.Do(func() {
		// Use simple_protocol to disable prepared statement cache at driver level
		// This prevents "cached plan must not change result type" errors during schema migrations
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC default_query_exec_mode=simple_protocol",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "pulse_secret"),
			getEnv("DB_NAME", "club_pulse"),
			getEnv("DB_PORT", "5432"),
		)

		var err error
		config := &gorm.Config{
			Logger:      logger.Default.LogMode(logger.Info),
			PrepareStmt: false, // Also disable GORM's prepared statement cache
		}

		// Retry logic for docker startup delay
		for i := 0; i < 5; i++ {
			DB, err = gorm.Open(postgres.Open(dsn), config)
			if err == nil {
				break
			}
			log.Printf("Failed to connect to database, retrying in 2 seconds... (%d/5)", i+1)
			time.Sleep(2 * time.Second)
		}

		if err != nil {
			log.Fatalf("Could not connect to database: %v", err)
		}

		sqlDB, err := DB.DB()
		if err != nil {
			log.Fatalf("Could not get database instance: %v", err)
		}

		// Connection pool settings
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		// Enable UUID extension
		if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
			log.Printf("Warning: Failed to create uuid-ossp extension: %v", err)
		}

		log.Println("Database connection established successfully")
	})
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	if DB == nil {
		InitDB()
	}
	return DB
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
