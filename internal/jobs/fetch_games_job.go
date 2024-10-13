package jobs

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
)

const delayBetweenRequests = 10 * time.Second

func transformToGame(appDetails *SteamAppDetails) domain.Game {
	var releaseDate time.Time
	var err error

	if appDetails.Data.ReleaseDate.Soon {
		releaseDate, err = time.Parse("January 2006", appDetails.Data.ReleaseDate.Date)
		if err != nil {
			log.Printf("Failed to parse release date: %v", err)
			releaseDate = time.Time{}
		}
	} else {
		releaseDate, err = time.Parse("2 Jan, 2006", appDetails.Data.ReleaseDate.Date)
		if err != nil {
			log.Printf("Failed to parse release date: %v", err)
			releaseDate = time.Time{}
		}
	}

	age, _ := strconv.Atoi(string(appDetails.Data.Age))

	return domain.Game{
		Slug:             utils.Slugify(appDetails.Data.Name),
		Title:            appDetails.Data.Name,
		About:            appDetails.Data.AboutTheGame,
		ShortDescription: appDetails.Data.ShortDescription,
		Description:      appDetails.Data.Description,
		Cover:            appDetails.Data.Background,
		Free:             appDetails.Data.IsFree,
		ReleaseDate:      releaseDate,
		Age:              age,
		Website:          &appDetails.Data.Website,
		Legal:            &appDetails.Data.Legal,
	}
}

func transformToDLC(appDetails *SteamAppDetails, gameID uint) domain.DLC {
	var releaseDate time.Time
	var err error

	if appDetails.Data.ReleaseDate.Soon {
		releaseDate, err = time.Parse("January 2006", appDetails.Data.ReleaseDate.Date)
		if err != nil {
			log.Printf("Failed to parse release date: %v", err)
			releaseDate = time.Time{}
		}
	} else {
		releaseDate, err = time.Parse("2 Jan, 2006", appDetails.Data.ReleaseDate.Date)
		if err != nil {
			log.Printf("Failed to parse release date: %v", err)
			releaseDate = time.Time{}
		}
	}

	dlcData := appDetails.Data

	return domain.DLC{
		Name:             dlcData.Name,
		About:            dlcData.AboutTheGame,
		ShortDescription: dlcData.ShortDescription,
		Description:      dlcData.Description,
		Cover:            dlcData.HeaderImage,
		Free:             dlcData.IsFree,
		ReleaseDate:      releaseDate,
		Legal:            &dlcData.Legal,
		GameID:           gameID,
	}
}

func PopulateSteamDatabaseJob(db *gorm.DB) {
	fmt.Println("Starting database population job...")

	apps, err := FetchSteamAppList()
	if err != nil {
		fmt.Printf("Failed to fetch app list: %v\n", err)
		return
	}

	for _, app := range apps {
		appDetails, err := FetchSteamAppDetails(app.AppID)
		if err != nil {
			fmt.Printf("Failed to fetch details for app ID %d: %v\n", app.AppID, err)
			continue
		}

		game := transformToGame(appDetails)
		if err := db.Create(&game).Error; err != nil {
			fmt.Printf("Failed to create game for app ID %d: %v\n", app.AppID, err)
			continue
		}

		gameID := game.ID

		requirements := map[string]any{
			"pc_requirements": map[string]string{
				"minimum":     appDetails.Data.PCRequirements.Minimum,
				"recommended": appDetails.Data.PCRequirements.Recommended,
			},
			"linux_requirements": map[string]string{
				"minimum":     appDetails.Data.LinuxRequirements.Minimum,
				"recommended": appDetails.Data.LinuxRequirements.Recommended,
			},
			"mac_requirements": map[string]string{
				"minimum":     appDetails.Data.MacRequirements.Minimum,
				"recommended": appDetails.Data.MacRequirements.Recommended,
			},
		}

		MapSteamGalleries(appDetails.Data.Screenshots, appDetails.Data.Movies, gameID, "games", db)
		MapSteamGamePrices(appDetails.Data.PriceOverview, gameID, strconv.Itoa(app.AppID), db)
		MapSteamSupport(appDetails.Data.Support, gameID, db)
		MapSteamGenresAndCategories(appDetails.Data.Genres, appDetails.Data.Categories, gameID, "games", db)
		MapSteamPublishersAndDevelopers(appDetails.Data.Developers, appDetails.Data.Publishers, gameID, "games", db)
		MapSteamSupportedLanguages(appDetails.Data.SupportedLanguages, gameID, "games", db)
		MapSteamRequirements(requirements, gameID, db)

		if len(appDetails.Data.DLC) > 0 {
			for _, dlcID := range appDetails.Data.DLC {
				appDlcDetails, err := FetchSteamAppDetails(dlcID)
				if err != nil {
					log.Printf("Failed to fetch details for app ID %d: %v\n", dlcID, err)
					continue
				}

				dlc := transformToDLC(appDlcDetails, gameID)
				if err := db.Create(&dlc).Error; err != nil {
					log.Printf("Failed to create dlc for app ID %d: %v\n", dlcID, err)
					continue
				}

				MapSteamDLCPrices(appDlcDetails.Data.PriceOverview, dlc.ID, strconv.Itoa(int(dlcID)), db)
				MapSteamGalleries(appDlcDetails.Data.Screenshots, appDlcDetails.Data.Movies, dlc.ID, "dlcs", db)
				MapSteamGenresAndCategories(appDlcDetails.Data.Genres, appDlcDetails.Data.Categories, dlc.ID, "dlcs", db)
				MapSteamPublishersAndDevelopers(appDlcDetails.Data.Developers, appDlcDetails.Data.Publishers, dlc.ID, "dlcs", db)
				MapSteamSupportedLanguages(appDlcDetails.Data.SupportedLanguages, dlc.ID, "dlcs", db)
			}
		}

		time.Sleep(delayBetweenRequests)
	}

	fmt.Println("Database population job completed.")
}

func FetchSteamOneByOneApp(db *gorm.DB, appID int) {
	fmt.Println("Starting database population job...")
	appDetails, err := FetchSteamAppDetails(appID)
	if err != nil {
		fmt.Printf("Failed to fetch details for app ID %d: %v\n", appID, err)
	}

	game := transformToGame(appDetails)
	if err := db.Create(&game).Error; err != nil {
		fmt.Printf("Failed to create game for app ID %d: %v\n", appID, err)
	}

	gameID := game.ID

	requirements := map[string]any{
		"pc_requirements": map[string]string{
			"minimum":     appDetails.Data.PCRequirements.Minimum,
			"recommended": appDetails.Data.PCRequirements.Recommended,
		},
		"linux_requirements": map[string]string{
			"minimum":     appDetails.Data.LinuxRequirements.Minimum,
			"recommended": appDetails.Data.LinuxRequirements.Recommended,
		},
		"mac_requirements": map[string]string{
			"minimum":     appDetails.Data.MacRequirements.Minimum,
			"recommended": appDetails.Data.MacRequirements.Recommended,
		},
	}

	MapSteamGalleries(appDetails.Data.Screenshots, appDetails.Data.Movies, gameID, "games", db)
	MapSteamGamePrices(appDetails.Data.PriceOverview, gameID, strconv.Itoa(appID), db)
	MapSteamSupport(appDetails.Data.Support, gameID, db)
	MapSteamGenresAndCategories(appDetails.Data.Genres, appDetails.Data.Categories, gameID, "games", db)
	MapSteamPublishersAndDevelopers(appDetails.Data.Developers, appDetails.Data.Publishers, gameID, "games", db)
	MapSteamSupportedLanguages(appDetails.Data.SupportedLanguages, gameID, "games", db)
	MapSteamRequirements(requirements, gameID, db)

	log.Printf("DLC length: %v", len(appDetails.Data.DLC))

	if len(appDetails.Data.DLC) > 0 {
		for _, dlcID := range appDetails.Data.DLC {
			appDlcDetails, err := FetchSteamAppDetails(dlcID)
			if err != nil {
				log.Printf("Failed to fetch details for app ID %d: %v\n", dlcID, err)
				continue
			}

			dlc := transformToDLC(appDlcDetails, gameID)
			if err := db.Create(&dlc).Error; err != nil {
				log.Printf("Failed to create dlc for app ID %d: %v\n", dlcID, err)
				continue
			}

			MapSteamDLCPrices(appDlcDetails.Data.PriceOverview, dlc.ID, strconv.Itoa(int(dlcID)), db)
			MapSteamGalleries(appDlcDetails.Data.Screenshots, appDlcDetails.Data.Movies, dlc.ID, "dlcs", db)
			MapSteamGenresAndCategories(appDlcDetails.Data.Genres, appDlcDetails.Data.Categories, dlc.ID, "dlcs", db)
			MapSteamPublishersAndDevelopers(appDlcDetails.Data.Developers, appDlcDetails.Data.Publishers, dlc.ID, "dlcs", db)
			MapSteamSupportedLanguages(appDlcDetails.Data.SupportedLanguages, dlc.ID, "dlcs", db)
		}
	}

	fmt.Println("Database population job completed.")
}
