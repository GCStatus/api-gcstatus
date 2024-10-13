package main

import (
	"flag"
	"fmt"
	"gcstatus/cmd/server/routes"
	"gcstatus/di"
	"gcstatus/internal/crons"
	"gcstatus/internal/jobs"
	"log"
	"os"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
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
		gameService,
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
		gameService,
		db,
	)

	c := cron.New()

	if _, err := c.AddFunc("@midnight", func() {
		crons.ResetMissions(db)
	}); err != nil {
		log.Fatalf("Failed to start cron: %+v", err)
	}

	// Register command to populate database
	populateSteamDBCmd := flag.Bool("populate-steam-db", false, "Populate the database with Steam games data")
	populateSteamOneDBCmd := flag.Bool("populate-steam-db-one", false, "Populate the database with only one Steam game data")
	appID := flag.Int("appID", 0, "App ID of the Steam game to populate (required if using populate-steam-db-one)")
	flag.Parse()

	// Start the cron scheduler
	c.Start()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)

	if *populateSteamDBCmd || *populateSteamOneDBCmd {
		if *populateSteamDBCmd {
			jobs.PopulateSteamDatabaseJob(db)
			fmt.Println("Database population job executed via command.")
		} else if *populateSteamOneDBCmd {
			if *appID == 0 {
				log.Fatalf("An appID is required when using -populate-steam-db-one")
			}
			jobs.FetchSteamOneByOneApp(db, *appID)
			fmt.Println("Database population for only one app job executed via command.")
		}
	} else {
		if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}

func BackgroundJobRunner(db *gorm.DB) {
	go jobs.PopulateSteamDatabaseJob(db)

	fmt.Println("Database population job started asynchronously.")
}
