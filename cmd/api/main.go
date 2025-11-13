package main

import (
	"fmt"
	"log"

	"golang-skeleton/configs"
	"golang-skeleton/database"
	"golang-skeleton/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load environment variables
	configs.LoadEnv()

	// Connect to database
	db, err := configs.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	log.Println("Database connected successfully!")

	// Auto migrate database
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed!")

	// Initialize Fiber app
	app := fiber.New()

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	port := configs.GetEnvOrDefault("APP_PORT", "8000")
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
