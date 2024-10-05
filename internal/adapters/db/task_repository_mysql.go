package db

import (
	"errors"
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"

	"gorm.io/gorm"
)

type TaskRepositoryMySQL struct {
	db *gorm.DB
}

func NewTaskRepositoryMySQL(db *gorm.DB) ports.TaskRepository {
	return &TaskRepositoryMySQL{db: db}
}

func (r *TaskRepositoryMySQL) GetTitleRequirementsByKey(actionKey string) ([]domain.TitleRequirement, error) {
	var requirements []domain.TitleRequirement
	err := r.db.Where("`key` = ?", actionKey).Find(&requirements).Error
	return requirements, err
}

func (r *TaskRepositoryMySQL) GetMissionRequirementsByKey(actionKey string) ([]domain.MissionRequirement, error) {
	var requirements []domain.MissionRequirement
	err := r.db.Where("`key` = ?", actionKey).Find(&requirements).Error
	return requirements, err
}

func (r *TaskRepositoryMySQL) GetOrCreateTitleProgress(userID uint, TitleRequirementID uint) (*domain.TitleProgress, error) {
	var progress domain.TitleProgress

	err := r.db.
		Joins("JOIN title_requirements ON title_requirements.id = title_progresses.title_requirement_id").
		Joins("JOIN titles ON titles.id = title_requirements.title_id").
		Where("title_requirements.id = ? AND title_progresses.user_id = ?", TitleRequirementID, userID).
		Where("titles.status NOT IN (?, ?)", domain.TitleUnavailable, domain.TitleCanceled).
		FirstOrCreate(&progress, domain.TitleProgress{
			UserID:             userID,
			TitleRequirementID: TitleRequirementID,
			Progress:           0,
			Completed:          false,
		}).Error

	return &progress, err
}

func (r *TaskRepositoryMySQL) GetOrCreateMissionProgress(userID uint, MissionRequirementID uint) (*domain.MissionProgress, error) {
	var progress domain.MissionProgress

	err := r.db.
		Joins("JOIN mission_requirements ON mission_requirements.id = mission_progresses.mission_requirement_id").
		Joins("JOIN missions ON missions.id = mission_requirements.mission_id").
		Where("mission_requirements.id = ? AND mission_progresses.user_id = ?", MissionRequirementID, userID).
		Where("missions.status NOT IN (?, ?)", domain.MissionUnavailable, domain.MissionCanceled).
		FirstOrCreate(&progress, domain.MissionProgress{
			UserID:               userID,
			MissionRequirementID: MissionRequirementID,
			Progress:             0,
			Completed:            false,
		}).Error

	return &progress, err
}

func (r *TaskRepositoryMySQL) UpdateTitleProgress(progress *domain.TitleProgress) error {
	return r.db.Save(progress).Error
}

func (r *TaskRepositoryMySQL) UpdateMissionProgress(progress *domain.MissionProgress) error {
	return r.db.Save(progress).Error
}

func (r *TaskRepositoryMySQL) UserHasTitle(userID uint, titleID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.UserTitle{}).
		Where("user_id = ? AND title_id = ?", userID, titleID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *TaskRepositoryMySQL) AwardTitleToUser(userID uint, titleID uint) error {
	tx := r.db.Begin()

	userTitle := domain.UserTitle{
		UserID:  userID,
		TitleID: titleID,
		Enabled: false,
	}

	if err := tx.Create(&userTitle).Error; err != nil {
		tx.Rollback()
		return err
	}

	var requirements []domain.TitleRequirement
	if err := tx.Where("title_id = ?", titleID).Find(&requirements).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, requirement := range requirements {
		var progress domain.TitleProgress

		err := tx.Where("user_id = ? AND title_requirement_id = ?", userID, requirement.ID).First(&progress).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			progress = domain.TitleProgress{
				UserID:             userID,
				TitleRequirementID: requirement.ID,
				Progress:           requirement.Goal,
				Completed:          true,
			}

			if err := tx.Create(&progress).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else if !progress.Completed {
			progress.Progress = requirement.Goal
			progress.Completed = true
			if err := tx.Save(&progress).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}
