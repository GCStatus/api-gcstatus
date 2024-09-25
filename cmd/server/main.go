package main

import (
	"fmt"
	"gcstatus/cmd/server/routes"
	"gcstatus/di"
	"log"
	"os"
)

func main() {
	// Initialize dependencies (repository, service, etc.)
	userService, authService, passwordResetService, levelService, profileService, _ := di.InitDependencies()

	// Setup routes with dependency injection
	r := routes.SetupRouter(
		userService,
		authService,
		passwordResetService,
		levelService,
		profileService,
	)

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Log the port before starting the server
	log.Printf("Starting server on port %s", port)

	// Start the server
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
