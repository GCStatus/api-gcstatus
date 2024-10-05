package main

import (
	"fmt"
	"gcstatus/cmd/server/routes"
	"gcstatus/di"
	"gcstatus/internal/crons"
	"log"
	"os"

	"github.com/robfig/cron/v3"
)

func main() {
	// Initialize dependencies (repository, service, etc.)
	userService,
		authService,
		passwordResetService,
		levelService,
		profileService,
		titleService,
		taskService,
		walletService,
		transactionService,
		notificationService,
		missionService,
		db := di.InitDependencies()

	// Setup routes with dependency injection
	r := routes.SetupRouter(
		userService,
		authService,
		passwordResetService,
		levelService,
		profileService,
		titleService,
		taskService,
		walletService,
		transactionService,
		notificationService,
		missionService,
	)

	c := cron.New()

	if _, err := c.AddFunc("@midnight", func() {
		crons.ResetMissions(db)
	}); err != nil {
		log.Fatalf("Failed to start cron: %+v", err)
	}

	// Start the cron scheduler
	c.Start()

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
