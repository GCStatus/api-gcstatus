package jobs

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"log"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

const (
	steamGamesAssociationsMorphsType = "games"
	steamDlcsAssociationsMorphsType  = "dlcs"
)

func MapSteamGalleries(screenshots []struct {
	Path string `json:"path_full"`
}, videos []struct {
	MP4 struct {
		Max string `json:"max"`
	} `json:"mp4"`
}, associateID uint, morphs string, db *gorm.DB) {
	if len(screenshots) > 0 {
		for _, screenshot := range screenshots {
			gallery := domain.Galleriable{
				S3:              false,
				Path:            screenshot.Path,
				MediaTypeID:     domain.PhotoTypeID,
				GalleriableID:   associateID,
				GalleriableType: morphs,
			}

			if err := db.Model(&domain.Galleriable{}).Create(&gallery).Error; err != nil {
				log.Printf("Failed to create an image on database for associateID: %v", associateID)
			}
		}
	}

	if len(videos) > 0 {
		for _, video := range videos {
			gallery := domain.Galleriable{
				S3:              false,
				Path:            video.MP4.Max,
				MediaTypeID:     domain.VideoTypeID,
				GalleriableID:   associateID,
				GalleriableType: morphs,
			}

			if err := db.Model(&domain.Galleriable{}).Create(&gallery).Error; err != nil {
				log.Printf("Failed to create a video on database for associateID: %v", associateID)
			}
		}
	}
}

func MapSteamGamePrices(priceOverview struct {
	Currency string `json:"currency"`
	Initial  uint   `json:"initial"`
	Final    uint   `json:"final"`
}, gameID uint, appID string, db *gorm.DB) {
	gameStore := domain.GameStore{
		Price:       priceOverview.Final,
		URL:         fmt.Sprintf("https://store.steampowered.com/app/%s", appID),
		GameID:      gameID,
		StoreID:     domain.SteamStoreID,
		StoreGameID: appID,
	}

	if err := db.Model(&domain.GameStore{}).Create(&gameStore).Error; err != nil {
		log.Printf("Failed to create a store for gameID: %v", gameID)
	}
}

func MapSteamDLCPrices(priceOverview struct {
	Currency string `json:"currency"`
	Initial  uint   `json:"initial"`
	Final    uint   `json:"final"`
}, dlcID uint, appID string, db *gorm.DB) {
	dlcStore := domain.DLCStore{
		Price:     priceOverview.Final,
		URL:       fmt.Sprintf("https://store.steampowered.com/app/%s", appID),
		DLCID:     dlcID,
		StoreID:   domain.SteamStoreID,
		StorDLCID: appID,
	}

	if err := db.Model(&domain.DLCStore{}).Create(&dlcStore).Error; err != nil {
		log.Printf("Failed to create a store for dlcID: %v", dlcID)
	}
}

func MapSteamGenresAndCategories(
	genres []struct {
		Name string `json:"description"`
	},
	categories []struct {
		Name string `json:"description"`
	},
	associateID uint,
	morphs string,
	db *gorm.DB,
) {
	if len(genres) > 0 {
		for _, genre := range genres {
			dbGenre := domain.Genre{
				Name: genre.Name,
				Slug: utils.Slugify(genre.Name),
			}

			err := db.Model(&domain.Genre{}).Where("name = ?", genre.Name).FirstOrCreate(&dbGenre).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Printf("Database error while fetching or creating genre: %v", err)
				continue
			}

			genreable := domain.Genreable{
				GenreableID:   associateID,
				GenreableType: morphs,
				GenreID:       dbGenre.ID,
			}
			if err := db.Create(&genreable).Error; err != nil {
				log.Printf("Failed to associate genre %s to game %v", genre.Name, associateID)
			}
		}
	}

	if len(categories) > 0 {
		for _, category := range categories {
			dbCategory := domain.Category{
				Name: category.Name,
				Slug: utils.Slugify(category.Name),
			}

			err := db.Model(&domain.Category{}).Where("name = ?", category.Name).FirstOrCreate(&dbCategory).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Printf("Database error while fetching or creating category: %v", err)
				continue
			}

			categoriable := domain.Categoriable{
				CategoriableID:   associateID,
				CategoriableType: morphs,
				CategoryID:       dbCategory.ID,
			}
			if err := db.Create(&categoriable).Error; err != nil {
				log.Printf("Failed to associate category %s to game %v", category.Name, associateID)
			}
		}
	}
}

func MapSteamPublishersAndDevelopers(
	developers []string,
	publishers []string,
	associateID uint,
	morphs string,
	db *gorm.DB,
) {
	if len(developers) > 0 {
		for _, developer := range developers {
			dbDeveloper := domain.Developer{
				Name:   developer,
				Acting: true,
			}

			err := db.Model(&domain.Developer{}).Where("name = ?", developer).FirstOrCreate(&dbDeveloper).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Printf("Database error while fetching or creating developer: %v", err)
				continue
			}

			if morphs == steamGamesAssociationsMorphsType {
				gameDeveloper := domain.GameDeveloper{
					GameID:      associateID,
					DeveloperID: dbDeveloper.ID,
				}
				if err := db.Create(&gameDeveloper).Error; err != nil {
					log.Printf("Failed to associate developer %s to game %v", developer, associateID)
				}
			} else if morphs == steamDlcsAssociationsMorphsType {
				dlcDeveloper := domain.DLCDeveloper{
					DLCID:       associateID,
					DeveloperID: dbDeveloper.ID,
				}
				if err := db.Create(&dlcDeveloper).Error; err != nil {
					log.Printf("Failed to associate developer %s to dlc %v", developer, associateID)
				}
			}
		}
	}

	if len(publishers) > 0 {
		for _, publisher := range publishers {
			dbPublisher := domain.Publisher{
				Name:   publisher,
				Acting: true,
			}

			err := db.Model(&domain.Publisher{}).Where("name = ?", publisher).FirstOrCreate(&dbPublisher).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Printf("Database error while fetching or creating publisher: %v", err)
				continue
			}

			if morphs == steamGamesAssociationsMorphsType {
				gamePublisher := domain.GamePublisher{
					GameID:      associateID,
					PublisherID: dbPublisher.ID,
				}
				if err := db.Create(&gamePublisher).Error; err != nil {
					log.Printf("Failed to associate publisher %s to game %v", publisher, associateID)
				}
			} else if morphs == steamDlcsAssociationsMorphsType {
				dlcPublisher := domain.DLCPublisher{
					DLCID:       associateID,
					PublisherID: dbPublisher.ID,
				}
				if err := db.Create(&dlcPublisher).Error; err != nil {
					log.Printf("Failed to associate publisher %s to dlc %v", publisher, associateID)
				}
			}
		}
	}
}

func MapSteamSupport(support struct {
	URL   string `json:"url"`
	Email string `json:"email"`
}, gameID uint, db *gorm.DB) {
	dbSupport := domain.GameSupport{
		URL:    &support.URL,
		Email:  &support.Email,
		GameID: gameID,
	}

	if err := db.Model(&domain.GameSupport{}).Create(&dbSupport).Error; err != nil {
		log.Fatalf("Failed to create a new support for game %v.", gameID)
	}
}

func MapSteamSupportedLanguages(supportedLanguages string, associateID uint, morphs string, db *gorm.DB) {
	langEntries := strings.Split(supportedLanguages, ",")
	audioSupportRegex := regexp.MustCompile(`(.*?)<strong>\*</strong>?`)

	for _, entry := range langEntries {
		entry = strings.TrimSpace(entry)
		entry = strings.ReplaceAll(entry, "<br>", "")
		match := audioSupportRegex.FindStringSubmatch(entry)

		langName := entry
		hasAudio := false

		if len(match) > 1 {
			langName = strings.TrimSpace(match[1])
			hasAudio = true
		}

		if langName == "" {
			continue
		}

		var language domain.Language
		if err := db.Where("name = ?", langName).FirstOrCreate(&language, domain.Language{Name: langName}).Error; err != nil {
			log.Printf("Failed while fetching or creating language %s: %v", langName, err)
			continue
		}

		if morphs == steamGamesAssociationsMorphsType {
			if err := db.Create(&domain.GameLanguage{
				GameID:     associateID,
				LanguageID: language.ID,
				Menu:       true,
				Dubs:       hasAudio,
				Subtitles:  true,
			}).Error; err != nil {
				log.Printf("Failed to associate language %s with game %v: %v", langName, associateID, err)
			}
		} else if morphs == steamDlcsAssociationsMorphsType {
			if err := db.Create(&domain.DLCLanguage{
				DLCID:      associateID,
				LanguageID: language.ID,
				Menu:       true,
				Dubs:       hasAudio,
				Subtitles:  true,
			}).Error; err != nil {
				log.Printf("Failed to associate language %s with dlc %v: %v", langName, associateID, err)
			}
		}
	}
}

func MapSteamRequirements(requirements map[string]any, associateID uint, db *gorm.DB) {
	requirementTypes := map[string]string{
		"pc_requirements":    "windows",
		"mac_requirements":   "mac",
		"linux_requirements": "linux",
	}

	for reqType, osType := range requirementTypes {
		reqMap, ok := requirements[reqType].(map[string]string)
		if !ok {
			log.Printf("No requirements found for %s", osType)
			continue
		}

		for potential, html := range reqMap {
			if isEmptyRequirement(html) {
				log.Printf("Skipping creation for empty requirement type %s and os type %s.", potential, osType)
				continue
			}

			os, dx, cpu, ram, gpu, storage, obs := extractRequirements(html)

			var reqType domain.RequirementType
			db.Where("potential = ? AND os = ?", potential, osType).FirstOrCreate(&reqType, domain.RequirementType{
				Potential: potential,
				OS:        osType,
			})

			if err := db.Create(&domain.Requirement{
				GameID:            associateID,
				RequirementTypeID: reqType.ID,
				OS:                os,
				DX:                dx,
				CPU:               cpu,
				RAM:               ram,
				GPU:               gpu,
				ROM:               storage,
				OBS:               obs,
				Network:           "N/A",
			}).Error; err != nil {
				log.Printf("Failed to create requirement for game %s", err.Error())
			}
		}
	}
}

func isEmptyRequirement(html string) bool {
	emptyPattern := regexp.MustCompile(`<strong>Minimum:</strong><br><ul class=\"bb_ul\"></ul>|<strong>Recommended:</strong><br><ul class=\"bb_ul\"></ul>`)
	return emptyPattern.MatchString(html)
}

func extractRequirements(html string) (string, string, string, string, string, string, *string) {
	var (
		osPattern      = regexp.MustCompile(`<strong>(?:OS|SO)\s*:</strong>\s*(.*?)\s*(?:<br|</li)`)
		dxPattern      = regexp.MustCompile(`<strong>DirectX:</strong>\s*(.*?)\s*(?:<br|</li)`)
		cpuPattern     = regexp.MustCompile(`<strong>Processor:</strong>\s*(.*?)\s*(?:<br|</li)`)
		ramPattern     = regexp.MustCompile(`<strong>Memory:</strong>\s*(.*?)\s*(?:<br|</li)`)
		gpuPattern     = regexp.MustCompile(`<strong>Graphics:</strong>\s*(.*?)\s*(?:<br|</li)`)
		storagePattern = regexp.MustCompile(`<strong>Storage:</strong>\s*(.*?)\s*(?:<br|</li)`)
		obsPattern     = regexp.MustCompile(`<strong>Additional Notes:</strong>\s*(.*?)\s*(?:<br|</li)`)
	)

	os := getStringMatch(osPattern, html)
	dx := getStringMatch(dxPattern, html)
	cpu := getStringMatch(cpuPattern, html)
	ram := getStringMatch(ramPattern, html)
	gpu := getStringMatch(gpuPattern, html)
	storage := getStringMatch(storagePattern, html)
	obs := getStringMatch(obsPattern, html)

	return os, dx, cpu, ram, gpu, storage, &obs
}

func getStringMatch(pattern *regexp.Regexp, html string) string {
	matches := pattern.FindStringSubmatch(html)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}
