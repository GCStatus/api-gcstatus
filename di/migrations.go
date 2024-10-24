package di

import (
	"gcstatus/internal/domain"
	"log"

	"gorm.io/gorm"
)

func MigrateModels(dbConn *gorm.DB) {
	models := []any{
		&domain.Reward{},
		&domain.Level{},
		&domain.Wallet{},
		&domain.User{},
		&domain.Profile{},
		&domain.PasswordReset{},
		&domain.Title{},
		&domain.TitleRequirement{},
		&domain.TitleProgress{},
		&domain.UserTitle{},
		&domain.TransactionType{},
		&domain.Transaction{},
		&domain.Notification{},
		&domain.Mission{},
		&domain.MissionRequirement{},
		&domain.MissionProgress{},
		&domain.UserMission{},
		&domain.UserMissionAssignment{},
		&domain.Genre{},
		&domain.Tag{},
		&domain.Platform{},
		&domain.Category{},
		&domain.Categoriable{},
		&domain.Genreable{},
		&domain.Taggable{},
		&domain.Platformable{},
		&domain.Language{},
		&domain.GameLanguage{},
		&domain.RequirementType{},
		&domain.Requirement{},
		&domain.Protection{},
		&domain.Cracker{},
		&domain.Crack{},
		&domain.TorrentProvider{},
		&domain.Torrent{},
		&domain.Publisher{},
		&domain.GamePublisher{},
		&domain.Developer{},
		&domain.GameDeveloper{},
		&domain.GameSupport{},
		&domain.Reviewable{},
		&domain.Viewable{},
		&domain.Heartable{},
		&domain.Critic{},
		&domain.Criticable{},
		&domain.Store{},
		&domain.GameStore{},
		&domain.Commentable{},
		&domain.Galleriable{},
		&domain.DLC{},
		&domain.DLCStore{},
		&domain.DLCPublisher{},
		&domain.DLCDeveloper{},
		&domain.DLCLanguage{},
		&domain.Game{},
		&domain.Banner{},
		&domain.Role{},
		&domain.Permission{},
		&domain.Roleable{},
		&domain.Permissionable{},
	}

	for _, model := range models {
		if err := dbConn.AutoMigrate(model); err != nil {
			log.Fatalf("Failed to migrate model %T: %v", model, err)
		}
	}
}
