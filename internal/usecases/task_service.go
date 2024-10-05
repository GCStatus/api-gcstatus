package usecases

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
)

type TaskService struct {
	repo ports.TaskRepository
}

func NewTaskService(repo ports.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) TrackTitleProgress(userID uint, actionKey string, increment int) error {
	requirements, err := s.GetTitleRequirementsByKey(actionKey)
	if err != nil {
		return err
	}

	for _, requirement := range requirements {
		progress, err := s.GetOrCreateTitleProgress(userID, requirement.ID)
		if err != nil {
			return err
		}

		if !progress.Completed {
			progress.Progress += increment
			if progress.Progress >= requirement.Goal {
				progress.Progress = requirement.Goal
				progress.Completed = true

				hasTitle, err := s.UserHasTitle(userID, requirement.TitleID)
				if err != nil {
					return err
				}

				if !hasTitle {
					err = s.AwardTitleToUser(userID, requirement.TitleID)
					if err != nil {
						return err
					}
				}
			}

			err = s.UpdateTitleProgress(progress)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *TaskService) TrackMissionProgress(userID uint, actionKey string, increment int) error {
	requirements, err := s.GetMissionRequirementsByKey(actionKey)
	if err != nil {
		return err
	}

	for _, requirement := range requirements {
		progress, err := s.GetOrCreateMissionProgress(userID, requirement.ID)
		if err != nil {
			return err
		}

		if !progress.Completed {
			progress.Progress += increment
			if progress.Progress >= requirement.Goal {
				progress.Progress = requirement.Goal
				progress.Completed = true
			}

			err = s.UpdateMissionProgress(progress)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *TaskService) GetTitleRequirementsByKey(actionKey string) ([]domain.TitleRequirement, error) {
	return s.repo.GetTitleRequirementsByKey(actionKey)
}

func (s *TaskService) GetMissionRequirementsByKey(actionKey string) ([]domain.MissionRequirement, error) {
	return s.repo.GetMissionRequirementsByKey(actionKey)
}

func (s *TaskService) GetOrCreateTitleProgress(userID uint, requirementID uint) (*domain.TitleProgress, error) {
	return s.repo.GetOrCreateTitleProgress(userID, requirementID)
}

func (s *TaskService) GetOrCreateMissionProgress(userID uint, requirementID uint) (*domain.MissionProgress, error) {
	return s.repo.GetOrCreateMissionProgress(userID, requirementID)
}

func (s *TaskService) UserHasTitle(userID uint, titleID uint) (bool, error) {
	return s.repo.UserHasTitle(userID, titleID)
}

func (s *TaskService) AwardTitleToUser(userID uint, titleID uint) error {
	return s.repo.AwardTitleToUser(userID, titleID)
}

func (s *TaskService) UpdateTitleProgress(progress *domain.TitleProgress) error {
	return s.repo.UpdateTitleProgress(progress)
}

func (s *TaskService) UpdateMissionProgress(progress *domain.MissionProgress) error {
	return s.repo.UpdateMissionProgress(progress)
}
