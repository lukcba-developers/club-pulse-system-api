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
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "postgres"),
			getEnv("DB_NAME", "club_pulse"),
			getEnv("DB_PORT", "5432"),
		)

		var err error
		config := &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
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
