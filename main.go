package main

import (
	"log"
	"os"

	"silic0n-wiki/config"
	"silic0n-wiki/database"
	"silic0n-wiki/routes"
)

func main() {
	if err := config.Load("config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := database.RunMigrations("./database/migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := os.MkdirAll(config.AppConfig.Media.UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	log.Printf("Server starting on port %d", config.AppConfig.Server.Port)
	routes.StartRouter()
}
