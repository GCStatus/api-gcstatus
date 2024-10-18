package tests

import (
	"errors"
	"fmt"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMissionRepositoryMySQL_FindById(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := testutils.Setup(t)
	repo := db.NewMissionRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		missionID       uint
		mockSetup       func()
		expectedError   error
		expectedMission *domain.Mission
	}{
		"success - mission found": {
			missionID: 1,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `missions` WHERE `missions`.`id` = ? AND `missions`.`deleted_at` IS NULL ORDER BY `missions`.`id` LIMIT ?")).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "mission", "description", "status", "for_all", "coins", "experience", "frequency", "reset_time", "created_at", "updated_at"}).
						AddRow(1, "Mission 1", "Description", "available", true, 10, 50, "daily", fixedTime, fixedTime, fixedTime))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rewards` WHERE `sourceable_type` = ? AND `rewards`.`sourceable_id` = ? AND `rewards`.`deleted_at` IS NULL")).
					WithArgs("missions", 1).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectedMission: &domain.Mission{
				ID:          1,
				Mission:     "Mission 1",
				Description: "Description",
				Status:      "available",
				ForAll:      true,
				Coins:       10,
				Experience:  50,
				Frequency:   "daily",
				ResetTime:   fixedTime,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				Rewards:     []domain.Reward{},
			},
			expectedError: nil,
		},
		"error - mission not found": {
			missionID: 2,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `missions` WHERE `missions`.`id` = ? AND `missions`.`deleted_at` IS NULL ORDER BY `missions`.`id` LIMIT ?")).
					WithArgs(2, 1).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectedMission: &domain.Mission{},
			expectedError:   gorm.ErrRecordNotFound,
		},
		"error - db failure": {
			missionID: 3,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `missions` WHERE `missions`.`id` = ? AND `missions`.`deleted_at` IS NULL ORDER BY `missions`.`id` LIMIT ?")).
					WithArgs(3, 1).
					WillReturnError(errors.New("db error"))
			},
			expectedMission: &domain.Mission{},
			expectedError:   errors.New("db error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()

			mission, err := repo.FindByID(tc.missionID)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedMission, mission)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMissionRepositoryMySQL_GetAllForUser(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := testutils.Setup(t)
	repo := db.NewMissionRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		userID           uint
		mockSetup        func()
		expectedError    error
		expectedMissions []*domain.Mission
	}{
		"success - missions found": {
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT `missions`.`id`,`missions`.`created_at`,`missions`.`updated_at`,`missions`.`deleted_at`,`missions`.`mission`,`missions`.`description`,`missions`.`status`,`missions`.`for_all`,`missions`.`coins`,`missions`.`experience`,`missions`.`frequency`,`missions`.`reset_time` FROM `missions` LEFT JOIN user_missions ON user_missions.mission_id = missions.id AND user_missions.user_id = ? WHERE (missions.for_all = ? OR user_missions.user_id = ?) AND missions.status NOT IN (?, ?) AND `missions`.`deleted_at` IS NULL")).
					WithArgs(1, true, 1, domain.MissionUnavailable, domain.MissionCanceled).
					WillReturnRows(sqlmock.NewRows([]string{"id", "mission", "description", "status", "for_all", "coins", "experience", "frequency", "reset_time", "created_at", "updated_at"}).
						AddRow(1, "Mission 1", "Description", "available", true, 10, 50, "daily", fixedTime, fixedTime, fixedTime))

				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `mission_requirements` WHERE `mission_requirements`.`mission_id` = ? AND `mission_requirements`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "mission_id"}).
						AddRow(1, 1))

				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `mission_progresses` WHERE `mission_progresses`.`mission_requirement_id` = ? AND user_id = ? AND `mission_progresses`.`deleted_at` IS NULL")).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "mission_requirement_id", "completed"}).
						AddRow(1, 1, 1, true))

				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `rewards` WHERE `sourceable_type` = ? AND `rewards`.`sourceable_id` = ? AND rewardable_type = ? AND `rewards`.`deleted_at` IS NULL")).
					WithArgs("missions", 1, "titles").
					WillReturnRows(sqlmock.NewRows([]string{"id", "sourceable_id", "sourceable_type", "rewardable_id", "rewardable_type"}).
						AddRow(1, 1, "missions", 1, "titles"))

				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `user_missions` WHERE `user_missions`.`mission_id` = ? AND user_id = ? AND `user_missions`.`deleted_at` IS NULL",
				)).WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "mission_id", "completed", "last_completed_at"}).
						AddRow(1, 1, 1, false, fixedTime))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `titles` WHERE id IN (?) AND `titles`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).
						AddRow(1, "Title 1"))
			},
			expectedMissions: []*domain.Mission{
				{
					ID:          1,
					Mission:     "Mission 1",
					Description: "Description",
					Status:      "available",
					ForAll:      true,
					Coins:       10,
					Experience:  50,
					Frequency:   "daily",
					ResetTime:   fixedTime,
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
					Rewards: []domain.Reward{
						{
							RewardableType: "titles",
							RewardableID:   1,
							Rewardable: &domain.Title{
								ID:    1,
								Title: "Title 1",
							},
						},
					},
					UserMission: []domain.UserMission{
						{
							ID:              1,
							Completed:       false,
							LastCompletedAt: fixedTime,
							UserID:          1,
							MissionID:       1,
						},
					},
				},
			},
			expectedError: nil,
		},
		"error - no missions found": {
			userID: 2,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT `missions`.`id`,`missions`.`created_at`,`missions`.`updated_at`,`missions`.`deleted_at`,`missions`.`mission`,`missions`.`description`,`missions`.`status`,`missions`.`for_all`,`missions`.`coins`,`missions`.`experience`,`missions`.`frequency`,`missions`.`reset_time` FROM `missions` LEFT JOIN user_missions ON user_missions.mission_id = missions.id AND user_missions.user_id = ? WHERE (missions.for_all = ? OR user_missions.user_id = ?) AND missions.status NOT IN (?, ?) AND `missions`.`deleted_at` IS NULL")).
					WithArgs(2, true, 2, domain.MissionUnavailable, domain.MissionCanceled).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectedMissions: []*domain.Mission{},
			expectedError:    nil,
		},
		"error - db failure": {
			userID: 3,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT `missions`.`id`,`missions`.`created_at`,`missions`.`updated_at`,`missions`.`deleted_at`,`missions`.`mission`,`missions`.`description`,`missions`.`status`,`missions`.`for_all`,`missions`.`coins`,`missions`.`experience`,`missions`.`frequency`,`missions`.`reset_time` FROM `missions` LEFT JOIN user_missions ON user_missions.mission_id = missions.id AND user_missions.user_id = ? WHERE (missions.for_all = ? OR user_missions.user_id = ?) AND missions.status NOT IN (?, ?) AND `missions`.`deleted_at` IS NULL")).
					WithArgs(3, true, 3, domain.MissionUnavailable, domain.MissionCanceled).
					WillReturnError(errors.New("db error"))
			},
			expectedMissions: nil,
			expectedError:    errors.New("db error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()

			missions, err := repo.GetAllForUser(tc.userID)

			assert.Equal(t, tc.expectedError, err)

			if len(missions) != len(tc.expectedMissions) {
				t.Fatalf("Expected %d missions, got %d", len(tc.expectedMissions), len(missions))
			}

			for i := range missions {
				if missions[i].ID != tc.expectedMissions[i].ID ||
					missions[i].Mission != tc.expectedMissions[i].Mission ||
					missions[i].Description != tc.expectedMissions[i].Description ||
					missions[i].Status != tc.expectedMissions[i].Status ||
					missions[i].ForAll != tc.expectedMissions[i].ForAll ||
					missions[i].Coins != tc.expectedMissions[i].Coins ||
					missions[i].Experience != tc.expectedMissions[i].Experience ||
					missions[i].Frequency != tc.expectedMissions[i].Frequency ||
					!missions[i].ResetTime.Equal(tc.expectedMissions[i].ResetTime) ||
					!missions[i].CreatedAt.Equal(tc.expectedMissions[i].CreatedAt) ||
					!missions[i].UpdatedAt.Equal(tc.expectedMissions[i].UpdatedAt) {
					t.Errorf("Expected mission %v, got %v", tc.expectedMissions[i], missions[i])
				}

				if len(missions[i].Rewards) != len(tc.expectedMissions[i].Rewards) {
					t.Fatalf("Expected %d rewards for mission %d, got %d", len(tc.expectedMissions[i].Rewards), missions[i].ID, len(missions[i].Rewards))
				}
				for j := range missions[i].Rewards {
					if missions[i].Rewards[j].RewardableType != tc.expectedMissions[i].Rewards[j].RewardableType ||
						missions[i].Rewards[j].RewardableID != tc.expectedMissions[i].Rewards[j].RewardableID {
						t.Errorf("Expected reward %v, got %v", tc.expectedMissions[i].Rewards[j], missions[i].Rewards[j])
					}

					if missions[i].Rewards[j].Rewardable != nil && tc.expectedMissions[i].Rewards[j].Rewardable != nil {
						title := missions[i].Rewards[j].Rewardable.(*domain.Title)
						expectedTitle := tc.expectedMissions[i].Rewards[j].Rewardable.(*domain.Title)
						if title.ID != expectedTitle.ID || title.Title != expectedTitle.Title {
							t.Errorf("Expected title %v, got %v", expectedTitle, title)
						}
					}
				}

				if len(missions[i].UserMission) != len(tc.expectedMissions[i].UserMission) {
					t.Fatalf("Expected %d user missions for mission %d, got %d", len(tc.expectedMissions[i].UserMission), missions[i].ID, len(missions[i].UserMission))
				}

				for j := range missions[i].UserMission {
					userMission := missions[i].UserMission[j]
					expectedUserMission := tc.expectedMissions[i].UserMission[j]

					if userMission.ID != expectedUserMission.ID ||
						userMission.Completed != expectedUserMission.Completed ||
						userMission.UserID != expectedUserMission.UserID ||
						userMission.MissionID != expectedUserMission.MissionID {
						t.Errorf("Expected user mission %v, got %v", expectedUserMission, userMission)
					}

					if userMission.LastCompletedAt != expectedUserMission.LastCompletedAt {
						t.Errorf("Expected LastCompletedAt %v, got %v", expectedUserMission.LastCompletedAt, userMission.LastCompletedAt)
					}
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMissionRepositoryMySQL_CompleteMission(t *testing.T) {
	gormDB, mock := testutils.Setup(t)
	repo := db.NewMissionRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		userID        uint
		missionID     uint
		mockSetup     func(userID uint, missionID uint)
		expectedError error
	}{
		"mission not found or unavailable": {
			userID:    1,
			missionID: 1,
			mockSetup: func(userID uint, missionID uint) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `missions` WHERE (id = ? AND status NOT IN (?, ?)) AND `missions`.`deleted_at` IS NULL ORDER BY `missions`.`id` LIMIT ?")).
					WithArgs(missionID, domain.MissionUnavailable, domain.MissionCanceled, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: fmt.Errorf("mission not found or unavailable: record not found"),
		},
		"user not assigned to this mission": {
			userID:    1,
			missionID: 2,
			mockSetup: func(userID uint, missionID uint) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `missions` WHERE (id = ? AND status NOT IN (?, ?)) AND `missions`.`deleted_at` IS NULL ORDER BY `missions`.`id` LIMIT ?")).
					WithArgs(missionID, domain.MissionUnavailable, domain.MissionCanceled, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "for_all"}).AddRow(missionID, false))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_mission_assignments` WHERE (user_id = ? AND mission_id = ?) AND `user_mission_assignments`.`deleted_at` IS NULL ORDER BY `user_mission_assignments`.`id` LIMIT ?")).
					WithArgs(userID, missionID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: fmt.Errorf("user is not assigned to this mission"),
		},
		"mission already completed by user": {
			userID:    1,
			missionID: 3,
			mockSetup: func(userID uint, missionID uint) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `missions` WHERE (id = ? AND status NOT IN (?, ?)) AND `missions`.`deleted_at` IS NULL ORDER BY `missions`.`id` LIMIT ?")).
					WithArgs(missionID, domain.MissionUnavailable, domain.MissionCanceled, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "for_all"}).AddRow(missionID, true))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_missions` WHERE (`user_missions`.`user_id` = ? AND `user_missions`.`mission_id` = ?) AND `user_missions`.`deleted_at` IS NULL ORDER BY `user_missions`.`id` LIMIT ?")).
					WithArgs(userID, missionID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "completed"}).AddRow(1, true))
			},
			expectedError: fmt.Errorf("mission already completed by user"),
		},
		"mission requirements not fully completed": {
			userID:    1,
			missionID: 4,
			mockSetup: func(userID uint, missionID uint) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `missions` WHERE (id = ? AND status NOT IN (?, ?)) AND `missions`.`deleted_at` IS NULL ORDER BY `missions`.`id` LIMIT ?")).
					WithArgs(missionID, domain.MissionUnavailable, domain.MissionCanceled, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "for_all"}).AddRow(missionID, true))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_missions` WHERE (`user_missions`.`user_id` = ? AND `user_missions`.`mission_id` = ?) AND `user_missions`.`deleted_at` IS NULL ORDER BY `user_missions`.`id` LIMIT ?")).
					WithArgs(userID, missionID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "completed"}).AddRow(1, false))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `mission_requirements` WHERE mission_id = ? AND `mission_requirements`.`deleted_at` IS NULL")).
					WithArgs(missionID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `mission_progresses` WHERE (mission_requirement_id = ? AND user_id = ?) AND `mission_progresses`.`deleted_at` IS NULL ORDER BY `mission_progresses`.`id` LIMIT ?")).
					WithArgs(1, userID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "completed"}).AddRow(1, false))
			},
			expectedError: fmt.Errorf("mission requirements not yet fully completed"),
		},
		"success - mission completed": {
			userID:    1,
			missionID: 5,
			mockSetup: func(userID uint, missionID uint) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `missions` WHERE (id = ? AND status NOT IN (?, ?)) AND `missions`.`deleted_at` IS NULL ORDER BY `missions`.`id` LIMIT ?")).
					WithArgs(missionID, domain.MissionUnavailable, domain.MissionCanceled, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "for_all"}).AddRow(missionID, true))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_missions` WHERE (`user_missions`.`user_id` = ? AND `user_missions`.`mission_id` = ?) AND `user_missions`.`deleted_at` IS NULL ORDER BY `user_missions`.`id` LIMIT ?")).
					WithArgs(userID, missionID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "completed"}).AddRow(1, false))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `mission_requirements` WHERE mission_id = ? AND `mission_requirements`.`deleted_at` IS NULL")).
					WithArgs(missionID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `mission_progresses` WHERE (mission_requirement_id = ? AND user_id = ?) AND `mission_progresses`.`deleted_at` IS NULL ORDER BY `mission_progresses`.`id` LIMIT ?")).
					WithArgs(1, userID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "completed"}).AddRow(1, true))

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_missions` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`completed`=?,`last_completed_at`=?,`user_id`=?,`mission_id`=? WHERE `user_missions`.`deleted_at` IS NULL AND `id` = ?")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), true, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup(tc.userID, tc.missionID)

			err := repo.CompleteMission(tc.userID, tc.missionID)

			if err == nil && tc.expectedError != nil {
				t.Errorf("Expected error %v, got nil", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Errorf("Expected nil, got error %v", err)
			} else if err != nil && tc.expectedError != nil {
				if err.Error() != tc.expectedError.Error() {
					t.Errorf("Expected error %v, got %v", tc.expectedError, err)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
