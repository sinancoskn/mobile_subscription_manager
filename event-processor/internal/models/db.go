package models

import (
	"event-processor/internal/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB initializes a new database connection
func NewDB(config *config.Config) (*gorm.DB, error) {
	// Configure GORM logging
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // Log writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // Disable colorful logs
		},
	)

	// Open the database connection
	db, err := gorm.Open(postgres.Open(config.DatabaseDSN), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	// Test the database connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Configure connection pool settings
	sqlDB.SetMaxIdleConns(10)                  // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)                 // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Maximum lifetime of a connection

	return db, nil
}
