package crons

import (
	"gcstatus/internal/domain"
	"log"
	"time"

	"gorm.io/gorm"
)

func ResetMissions(db *gorm.DB) {
	log.Printf("start running cron...")

	now := time.Now()

	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var missions []domain.Mission

	db.Where("reset_time < ? AND status NOT IN (?,?) AND frequency != ?", now, domain.MissionCanceled, domain.MissionUnavailable, domain.OneTimeMission).Find(&missions)

	for _, mission := range missions {
		if mission.ForAll {
			resetAllUserProgressForMission(db, mission)
		} else {
			resetSpecificUserProgressForMission(db, mission)
		}

		switch mission.Frequency {
		case "daily":
			mission.ResetTime = midnight.Add(24 * time.Hour)
		case "weekly":
			mission.ResetTime = midnight.AddDate(0, 0, 7)
		case "monthly":
			mission.ResetTime = midnight.AddDate(0, 1, 0)
		}

		db.Save(&mission)
	}

	log.Printf("cron runned successfully!")
}

func resetAllUserProgressForMission(db *gorm.DB, mission domain.Mission) {
	var users []domain.User

	db.Find(&users)

	for _, user := range users {
		resetProgressForUser(db, user.ID, mission)
	}
}

func resetSpecificUserProgressForMission(db *gorm.DB, mission domain.Mission) {
	var userMissions []domain.UserMission

	db.Where("mission_id = ?", mission.ID).Find(&userMissions)

	for _, userMission := range userMissions {
		resetProgressForUser(db, userMission.UserID, mission)
	}
}

func resetProgressForUser(db *gorm.DB, userID uint, mission domain.Mission) {
	var missionProgress []domain.MissionProgress

	db.Where("mission_requirement_id IN (?) AND user_id = ?", db.Table("mission_requirements").Select("id").Where("mission_id = ?", mission.ID), userID).Find(&missionProgress)

	if len(missionProgress) == 0 {
		return
	}

	for _, progress := range missionProgress {
		progress.Progress = 0
		progress.Completed = false
		db.Save(&progress)
	}

	var userMission domain.UserMission
	if err := db.Where("mission_id = ? AND user_id = ?", mission.ID, userID).First(&userMission).Error; err == nil {
		userMission.Completed = false
		userMission.LastCompletedAt = time.Time{}
		db.Save(&userMission)
	}
}
