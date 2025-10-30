package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewMasterDB creates the singleton master DB connection (pool configured)
func NewMasterDB() (*gorm.DB, error) {
	// read DSN from env for production
	dsn := os.Getenv("MASTER_DSN")
	if dsn == "" {
		// example local dsn
		dsn = "host=localhost user=AAAAAAA password=AAAAA@123 dbname=masterdb port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect master db: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql DB: %w", err)
	}
	// Master pool config (tune as needed)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(60 * 60 * 1e9) // 1 hour in ns
	return db, nil
}
