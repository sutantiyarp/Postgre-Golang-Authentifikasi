// package config

// import (
// 	// "database/sql"
// 	"hello-fiber/route"
// 	"hello-fiber/middleware"
// 	"github.com/gofiber/fiber/v2"
// 	"hello-fiber/database"  // Mengimpor package database
// )

// func NewApp() *fiber.App {
// 	// Connect to the database
// 	db := database.ConnectDB()

// 	// Initialize the Fiber application
// 	app := fiber.New()

// 	// Middleware
// 	app.Use(middleware.LoggerMiddleware)

// 	// Set up routes, passing db as a dependency to the route handler
// 	route.SetupRoutes(app, db)

// 	return app
// }

package config

import (
	// "database/sql"
	"hello-fiber/route"
	"hello-fiber/middleware"
	"github.com/gofiber/fiber/v2"
	"hello-fiber/database"  // Mengimpor package database
)

func NewApp() *fiber.App {
	mongoClient := database.ConnectMongoDB()
	_ = mongoClient // Simpan reference jika diperlukan

	// Initialize the Fiber application
	app := fiber.New(fiber.Config{
		BodyLimit: 2 * 1024 * 1024, // Set body limit to 2MB for file uploads
	})

	// Middleware
	app.Use(middleware.LoggerMiddleware)

	app.Static("/file", "./uploads")

	// Set up routes, passing db as a dependency to the route handler
	route.SetupRoutes(app, nil) // Pass nil untuk db karena sudah global di database.MongoDB

	return app
}