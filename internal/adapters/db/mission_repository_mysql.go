package db

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
	"time"

	"gorm.io/gorm"
)

type MissionRepositoryMySQL struct {
	db *gorm.DB
}

func NewMissionRepositoryMySQL(db *gorm.DB) ports.MissionRepository {
	return &MissionRepositoryMySQL{db: db}
}

func (h *MissionRepositoryMySQL) FindByID(missionID uint) (*domain.Mission, error) {
	var mission *domain.Mission
	err := h.db.Preload("Rewards").First(&mission, missionID).Error
	return mission, err
}

func (h *MissionRepositoryMySQL) GetAllForUser(userID uint) ([]*domain.Mission, error) {
	var missions []*domain.Mission

	err := h.db.
		Model(&domain.Mission{}).
		Joins("LEFT JOIN user_missions ON user_missions.mission_id = missions.id AND user_missions.user_id = ?", userID).
		Where("missions.for_all = ? OR user_missions.user_id = ?", true, userID).
		Where("missions.status NOT IN (?, ?)", domain.MissionUnavailable, domain.MissionCanceled).
		Preload("MissionRequirements").
		Preload("MissionRequirements.MissionProgress", "user_id = ?", userID).
		Preload("UserMission", "user_id = ?", userID).
		Preload("Rewards", "rewardable_type = ?", "titles").
		Find(&missions).
		Error
	if err != nil {
		return nil, err
	}

	var titleIDs []uint
	for _, mission := range missions {
		for _, reward := range mission.Rewards {
			if reward.RewardableType == domain.RewardableTypeTitles {
				titleIDs = append(titleIDs, reward.RewardableID)
			}
		}
	}

	var titles []domain.Title
	if len(titleIDs) > 0 {
		err = h.db.Where("id IN (?)", titleIDs).Find(&titles).Error
		if err != nil {
			return nil, err
		}
	}

	titleMap := make(map[uint]domain.Title)
	for _, title := range titles {
		titleMap[title.ID] = title
	}

	for _, mission := range missions {
		for i := range mission.Rewards {
			reward := &mission.Rewards[i]
			if reward.RewardableType == domain.RewardableTypeTitles {
				if title, ok := titleMap[reward.RewardableID]; ok {
					reward.Rewardable = &title
				}
			}
		}
	}

	return missions, nil
}

func (h *MissionRepositoryMySQL) CompleteMission(userID uint, missionID uint) error {
	var mission domain.Mission
	if err := h.db.Where("id = ? AND status NOT IN (?, ?)", missionID, domain.MissionUnavailable, domain.MissionCanceled).First(&mission).Error; err != nil {
		return fmt.Errorf("mission not found or unavailable: %w", err)
	}

	if !mission.ForAll {
		var userMissionAssignment domain.UserMissionAssignment
		if err := h.db.Where("user_id = ? AND mission_id = ?", userID, missionID).First(&userMissionAssignment).Error; err != nil {
			return fmt.Errorf("user is not assigned to this mission")
		}
	}

	var userMission domain.UserMission
	if err := h.db.FirstOrCreate(&userMission, domain.UserMission{
		UserID:    userID,
		MissionID: missionID,
	}).Error; err != nil {
		return fmt.Errorf("error creating or finding user mission entry: %w", err)
	}

	if userMission.Completed {
		return fmt.Errorf("mission already completed by user")
	}

	var requirements []domain.MissionRequirement
	if err := h.db.Where("mission_id = ?", missionID).Find(&requirements).Error; err != nil {
		return fmt.Errorf("error fetching mission requirements: %w", err)
	}

	for _, req := range requirements {
		var progress domain.MissionProgress
		if err := h.db.Where("mission_requirement_id = ? AND user_id = ?", req.ID, userID).First(&progress).Error; err != nil {
			return fmt.Errorf("error fetching user progress: %w", err)
		}

		if !progress.Completed {
			return fmt.Errorf("mission requirements not yet fully completed")
		}
	}

	userMission.Completed = true
	userMission.LastCompletedAt = time.Now()

	if err := h.db.Save(&userMission).Error; err != nil {
		return fmt.Errorf("error updating user mission completion: %w", err)
	}

	return nil
}
