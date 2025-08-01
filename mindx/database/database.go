package database

import (
	"fmt"

	"mindx/config"
	"mindx/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes database connection and performs migrations
func InitDB(config config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate models - first RiskEvaluation, then Student to avoid circular references
	err = db.AutoMigrate(&models.RiskEvaluation{})
	if err != nil {
		return nil, err
	}
	
	err = db.AutoMigrate(&models.Student{})
	if err != nil {
		return nil, err
	}

	return db, nil
}