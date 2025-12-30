package main

import (
	"log"
	"os"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
)

func main() {
	log.Println("Starting explicit migration...")

	// Set envs for local (override if needed)
	os.Setenv("DB_PASSWORD", "pulse_secret")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_NAME", "club_pulse")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")

	database.InitDB()
	db := database.GetDB()

	log.Println("Migrating MaintenanceTaskModel...")
	err := db.AutoMigrate(&repository.MaintenanceTaskModel{})
	if err != nil {
		log.Fatalf("MaintenanceTaskModel migration failed: %v", err)
	}

	log.Println("Migrating EquipmentModel...")
	err = db.AutoMigrate(&repository.EquipmentModel{})
	if err != nil {
		log.Fatalf("EquipmentModel migration failed: %v", err)
	}

	log.Println("Migration successful!")
}
