package main

import (
	"hello-fiber/config"
	"log"
	// "os"
	// "database/sql"
	"github.com/joho/godotenv"
	// "github.com/lib/pq"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create and start the Fiber app
	app := config.NewApp()

	// Start the server
	log.Fatal(app.Listen(":3000"))
}