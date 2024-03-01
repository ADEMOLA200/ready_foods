package main

import (
	"github.com/ADEMOLA200/danas-food/database"
	"github.com/ADEMOLA200/danas-food/middleware"
	"github.com/ADEMOLA200/danas-food/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Connect to the database
	database.ConnectDB()
	
    // Create a new Fiber app
    app := fiber.New()

	// Register middleware
	app.Use(middleware.LoggingMiddleware())
	app.Use(middleware.ErrorHandlingMiddleware())
	app.Use(middleware.CorsMiddleware())
	//app.Use(middleware.IsAuthenticated)

	// Initialize routes
	routes.SetupRoutes(app)

	// Start the server on port 3000
	app.Listen(":3000")
}
