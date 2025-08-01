package main

import (
	"log"

	"mindx/config"
	"mindx/database"
	"mindx/router"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize router
	r := router.InitRouter(db)

	// Start server
	log.Printf("Server starting on %s", cfg.Server.Address)
	if err := r.Start(cfg.Server.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}