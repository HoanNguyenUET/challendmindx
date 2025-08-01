package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Risk     RiskConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Address string
}

// RiskConfig holds risk evaluation configuration
type RiskConfig struct {
	AttendanceThreshold float64
	AssignmentThreshold float64
	ContactThreshold    int
	LowRiskThreshold    int
	MediumRiskThreshold int
	HighRiskThreshold   int
}

// LoadConfig loads configuration from environment variables
// with sensible defaults
func LoadConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "studentrisk"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Address: getEnv("SERVER_ADDRESS", ":8080"),
		},
		Risk: RiskConfig{
			AttendanceThreshold: getEnvFloat("RISK_ATTENDANCE_THRESHOLD", 75.0),
			AssignmentThreshold: getEnvFloat("RISK_ASSIGNMENT_THRESHOLD", 50.0),
			ContactThreshold:    getEnvInt("RISK_CONTACT_THRESHOLD", 2),
			LowRiskThreshold:    getEnvInt("RISK_LOW_THRESHOLD", 0),
			MediumRiskThreshold: getEnvInt("RISK_MEDIUM_THRESHOLD", 2),
			HighRiskThreshold:   getEnvInt("RISK_HIGH_THRESHOLD", 3),
		},
	}
}

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get environment variable as float with default value
func getEnvFloat(key string, defaultValue float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		var result float64
		_, err := fmt.Sscanf(value, "%f", &result)
		if err == nil {
			return result
		}
	}
	return defaultValue
}

// Helper function to get environment variable as int with default value
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		var result int
		_, err := fmt.Sscanf(value, "%d", &result)
		if err == nil {
			return result
		}
	}
	return defaultValue
}