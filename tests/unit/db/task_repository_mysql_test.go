package tests

import (
	"errors"
	"fmt"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	"gcstatus/tests"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTaskRepositoryMySQL_GetTitleRequirementsByKey(t *testing.T) {
	gormDB, mock := tests.Setup(t)

	repo := db.NewTaskRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		key        string
		expectReqs []domain.TitleRequirement
		expectErr  bool
		expectMsg  string
	}{
		"valid requirements key": {
			key: "valid_key",
			expectReqs: []domain.TitleRequirement{
				{
					ID:          1,
					Task:        "Task 1",
					Description: "Task 1",
					Key:         "valid_key",
					Goal:        10,
				},
				{
					ID:          2,
					Task:        "Task 2",
					Description: "Task 2",
					Key:         "valid_key",
					Goal:        10,
				},
			},
			expectErr: false,
		},
		"no requirements found": {
			key:        "invalid_key",
			expectReqs: []domain.TitleRequirement{},
			expectErr:  false,
		},
		"error - db failure": {
			key:        "error_key",
			expectReqs: nil,
			expectErr:  true,
			expectMsg:  "db error",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			switch name {
			case "valid requirements key":
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `title_requirements` WHERE `key` = ?")).
					WithArgs(tc.key).
					WillReturnRows(sqlmock.NewRows([]string{"id", "task", "description", "key", "goal"}).
						AddRow(1, "Task 1", "Task 1", "valid_key", 10).
						AddRow(2, "Task 2", "Task 2", "valid_key", 10))

			case "no requirements found":
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `title_requirements` WHERE `key` = ?")).
					WithArgs(tc.key).
					WillReturnRows(sqlmock.NewRows([]string{}))

			case "error - db failure":
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `title_requirements` WHERE `key` = ?")).
					WithArgs(tc.key).
					WillReturnError(errors.New("db error"))
			}

			reqs, err := repo.GetTitleRequirementsByKey(tc.key)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectReqs, reqs)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskRepositoryMySQL_GetMissionRequirementsByKey(t *testing.T) {
	gormDB, mock := tests.Setup(t)

	repo := db.NewTaskRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		key        string
		expectReqs []domain.MissionRequirement
		expectErr  bool
		expectMsg  string
	}{
		"valid requirements key": {
			key: "valid_key",
			expectReqs: []domain.MissionRequirement{
				{
					ID:          1,
					Task:        "Task 1",
					Description: "Task 1",
					Key:         "valid_key",
					Goal:        10,
				},
				{
					ID:          2,
					Task:        "Task 2",
					Description: "Task 2",
					Key:         "valid_key",
					Goal:        10,
				},
			},
			expectErr: false,
		},
		"no requirements found": {
			key:        "invalid_key",
			expectReqs: []domain.MissionRequirement{},
			expectErr:  false,
		},
		"error - db failure": {
			key:        "error_key",
			expectReqs: nil,
			expectErr:  true,
			expectMsg:  "db error",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			switch name {
			case "valid requirements key":
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `mission_requirements` WHERE `key` = ?")).
					WithArgs(tc.key).
					WillReturnRows(sqlmock.NewRows([]string{"id", "task", "description", "key", "goal"}).
						AddRow(1, "Task 1", "Task 1", "valid_key", 10).
						AddRow(2, "Task 2", "Task 2", "valid_key", 10))

			case "no requirements found":
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `mission_requirements` WHERE `key` = ?")).
					WithArgs(tc.key).
					WillReturnRows(sqlmock.NewRows([]string{}))

			case "error - db failure":
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `mission_requirements` WHERE `key` = ?")).
					WithArgs(tc.key).
					WillReturnError(errors.New("db error"))
			}

			reqs, err := repo.GetMissionRequirementsByKey(tc.key)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectReqs, reqs)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskRepositoryMySQL_UpdateTitleProgress(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		titleProgressID     uint
		expectTitleProgress domain.TitleProgress
		mock                func(mock sqlmock.Sqlmock, progress domain.TitleProgress)
		expectErr           bool
		expectMsg           string
	}{
		"valid progress": {
			titleProgressID: 1,
			expectTitleProgress: domain.TitleProgress{
				ID:                 1,
				Progress:           5,
				Completed:          false,
				CreatedAt:          fixedTime,
				UpdatedAt:          fixedTime,
				UserID:             1,
				TitleRequirementID: 1,
			},
			mock: func(mock sqlmock.Sqlmock, progress domain.TitleProgress) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `title_progresses` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`progress`=?,`completed`=?,`user_id`=?,`title_requirement_id`=? WHERE `title_progresses`.`deleted_at` IS NULL AND `id` = ?")).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						progress.Progress,
						progress.Completed,
						progress.UserID,
						progress.TitleRequirementID,
						progress.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		"can create progress if not found": {
			titleProgressID:     999,
			expectTitleProgress: domain.TitleProgress{},
			mock: func(mock sqlmock.Sqlmock, progress domain.TitleProgress) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `title_progresses` (`created_at`,`updated_at`,`deleted_at`,`progress`,`completed`,`user_id`,`title_requirement_id`) VALUES (?,?,?,?,?,?,?)")).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						progress.Progress,
						progress.Completed,
						progress.UserID,
						progress.TitleRequirementID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		"error - db failure": {
			titleProgressID:     1,
			expectTitleProgress: domain.TitleProgress{},
			mock: func(mock sqlmock.Sqlmock, progress domain.TitleProgress) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `title_progresses` (`created_at`,`updated_at`,`deleted_at`,`progress`,`completed`,`user_id`,`title_requirement_id`) VALUES (?,?,?,?,?,?,?)")).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						progress.Progress,
						progress.Completed,
						progress.UserID,
						progress.TitleRequirementID,
					).
					WillReturnError(errors.New("db failure"))
				mock.ExpectRollback()
			},
			expectErr: true,
			expectMsg: "db failure",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := tests.Setup(t)

			repo := db.NewTaskRepositoryMySQL(gormDB)
			tc.mock(mock, tc.expectTitleProgress)

			err := repo.UpdateTitleProgress(&tc.expectTitleProgress)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectTitleProgress, tc.expectTitleProgress)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestTaskRepositoryMySQL_UpdateMissionProgress(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		missionProgressID     uint
		expectMissionProgress domain.MissionProgress
		mock                  func(mock sqlmock.Sqlmock, progress domain.MissionProgress)
		expectErr             bool
		expectMsg             string
	}{
		"valid progress": {
			missionProgressID: 1,
			expectMissionProgress: domain.MissionProgress{
				ID:                   1,
				Progress:             5,
				Completed:            false,
				CreatedAt:            fixedTime,
				UpdatedAt:            fixedTime,
				UserID:               1,
				MissionRequirementID: 1,
			},
			mock: func(mock sqlmock.Sqlmock, progress domain.MissionProgress) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `mission_progresses` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`progress`=?,`completed`=?,`user_id`=?,`mission_requirement_id`=? WHERE `mission_progresses`.`deleted_at` IS NULL AND `id` = ?")).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						progress.Progress,
						progress.Completed,
						progress.UserID,
						progress.MissionRequirementID,
						progress.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		"can create progress if not found": {
			missionProgressID:     999,
			expectMissionProgress: domain.MissionProgress{},
			mock: func(mock sqlmock.Sqlmock, progress domain.MissionProgress) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `mission_progresses` (`created_at`,`updated_at`,`deleted_at`,`progress`,`completed`,`user_id`,`mission_requirement_id`) VALUES (?,?,?,?,?,?,?)")).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						progress.Progress,
						progress.Completed,
						progress.UserID,
						progress.MissionRequirementID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		"error - db failure": {
			missionProgressID:     1,
			expectMissionProgress: domain.MissionProgress{},
			mock: func(mock sqlmock.Sqlmock, progress domain.MissionProgress) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `mission_progresses` (`created_at`,`updated_at`,`deleted_at`,`progress`,`completed`,`user_id`,`mission_requirement_id`) VALUES (?,?,?,?,?,?,?)")).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						progress.Progress,
						progress.Completed,
						progress.UserID,
						progress.MissionRequirementID,
					).
					WillReturnError(errors.New("db failure"))
				mock.ExpectRollback()
			},
			expectErr: true,
			expectMsg: "db failure",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := tests.Setup(t)

			repo := db.NewTaskRepositoryMySQL(gormDB)
			tc.mock(mock, tc.expectMissionProgress)

			err := repo.UpdateMissionProgress(&tc.expectMissionProgress)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectMissionProgress, tc.expectMissionProgress)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestTaskRepositoryMySQL_UserHasTitle(t *testing.T) {
	gormDB, mock := tests.Setup(t)
	r := db.NewTaskRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		userID    uint
		titleID   uint
		mock      func()
		expectErr bool
		expectMsg string
	}{
		"user has title": {
			userID:  1,
			titleID: 1,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `user_titles` WHERE (user_id = ? AND title_id = ?) AND `user_titles`.`deleted_at` IS NULL")).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectErr: false,
		},
		"user has not title": {
			userID:  1,
			titleID: 2,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `user_titles` WHERE (user_id = ? AND title_id = ?) AND `user_titles`.`deleted_at` IS NULL")).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectErr: true,
			expectMsg: "user has not title",
		},
		"user not found": {
			userID:  99,
			titleID: 1,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `user_titles` WHERE (user_id = ? AND title_id = ?) AND `user_titles`.`deleted_at` IS NULL")).
					WithArgs(99, 1).
					WillReturnError(fmt.Errorf("user not found"))
			},
			expectErr: true,
			expectMsg: "user not found",
		},
		"title not found": {
			userID:  1,
			titleID: 99,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `user_titles` WHERE (user_id = ? AND title_id = ?) AND `user_titles`.`deleted_at` IS NULL")).
					WithArgs(1, 99).
					WillReturnError(fmt.Errorf("title not found"))
			},
			expectErr: true,
			expectMsg: "title not found",
		},
		"db failure": {
			userID:  1,
			titleID: 1,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `user_titles` WHERE (user_id = ? AND title_id = ?) AND `user_titles`.`deleted_at` IS NULL")).
					WithArgs(1, 1).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectErr: true,
			expectMsg: "database error",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mock()

			result, err := r.UserHasTitle(tc.userID, tc.titleID)

			if tc.expectErr && err != nil {
				assert.Equal(t, err.Error(), tc.expectMsg)
			}

			if !tc.expectErr && !result {
				t.Error("expected true, got false")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestTaskRepositoryMySQL_GetOrCreateTitleProgress(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := tests.Setup(t)
	r := db.NewTaskRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		userID             uint
		titleRequirementID uint
		mock               func()
		expectErr          bool
		expectProgress     *domain.TitleProgress
	}{
		"existing progress": {
			userID:             1,
			titleRequirementID: 1,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT `title_progresses`.`id`,`title_progresses`.`created_at`,`title_progresses`.`updated_at`,`title_progresses`.`deleted_at`,`title_progresses`.`progress`,`title_progresses`.`completed`,`title_progresses`.`user_id`,`title_progresses`.`title_requirement_id` FROM `title_progresses` JOIN title_requirements ON title_requirements.id = title_progresses.title_requirement_id JOIN titles ON titles.id = title_requirements.title_id WHERE (title_requirements.id = ? AND title_progresses.user_id = ?) AND titles.status NOT IN (?, ?) AND (`title_progresses`.`user_id` = ? AND `title_progresses`.`title_requirement_id` = ?) AND `title_progresses`.`deleted_at` IS NULL ORDER BY `title_progresses`.`id` LIMIT ?")).
					WithArgs(1, 1, domain.TitleUnavailable, domain.TitleCanceled, 1, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title_requirement_id", "progress", "completed", "created_at", "updated_at"}).
						AddRow(1, 1, 1, 0, false, fixedTime, fixedTime))
			},
			expectErr: false,
			expectProgress: &domain.TitleProgress{
				ID:                 1,
				Progress:           0,
				Completed:          false,
				UserID:             1,
				TitleRequirementID: 1,
			},
		},
		"create new progress": {
			userID:             1,
			titleRequirementID: 2,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT `title_progresses`.`id`,`title_progresses`.`created_at`,`title_progresses`.`updated_at`,`title_progresses`.`deleted_at`,`title_progresses`.`progress`,`title_progresses`.`completed`,`title_progresses`.`user_id`,`title_progresses`.`title_requirement_id` FROM `title_progresses` JOIN title_requirements ON title_requirements.id = title_progresses.title_requirement_id JOIN titles ON titles.id = title_requirements.title_id WHERE (title_requirements.id = ? AND title_progresses.user_id = ?) AND titles.status NOT IN (?, ?) AND (`title_progresses`.`user_id` = ? AND `title_progresses`.`title_requirement_id` = ?) AND `title_progresses`.`deleted_at` IS NULL ORDER BY `title_progresses`.`id` LIMIT ?")).
					WithArgs(2, 1, domain.TitleUnavailable, domain.TitleCanceled, 1, 2, 1).
					WillReturnRows(sqlmock.NewRows([]string{}))

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `title_progresses` (`created_at`,`updated_at`,`deleted_at`,`progress`,`completed`,`user_id`,`title_requirement_id`) VALUES (?,?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 0, false, 1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
			expectProgress: &domain.TitleProgress{
				ID:                 1,
				Progress:           0,
				Completed:          false,
				UserID:             1,
				TitleRequirementID: 2,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mock()

			progress, err := r.GetOrCreateTitleProgress(tc.userID, tc.titleRequirementID)

			if (err != nil) != tc.expectErr {
				t.Errorf("expected error: %v, got: %v", tc.expectErr, err)
			} else {
				assert.Equal(t, progress.ID, tc.expectProgress.ID)
				assert.Equal(t, progress.Progress, tc.expectProgress.Progress)
				assert.Equal(t, progress.Completed, tc.expectProgress.Completed)
				assert.Equal(t, progress.UserID, tc.expectProgress.UserID)
				assert.Equal(t, progress.TitleRequirementID, tc.expectProgress.TitleRequirementID)

				assert.NotZero(t, progress.CreatedAt)
				assert.NotZero(t, progress.UpdatedAt)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestTaskRepositoryMySQL_GetOrCreateMissionProgress(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := tests.Setup(t)
	r := db.NewTaskRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		userID               uint
		missionRequirementID uint
		mock                 func()
		expectErr            bool
		expectProgress       *domain.MissionProgress
	}{
		"existing progress": {
			userID:               1,
			missionRequirementID: 1,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT `mission_progresses`.`id`,`mission_progresses`.`created_at`,`mission_progresses`.`updated_at`,`mission_progresses`.`deleted_at`,`mission_progresses`.`progress`,`mission_progresses`.`completed`,`mission_progresses`.`user_id`,`mission_progresses`.`mission_requirement_id` FROM `mission_progresses` JOIN mission_requirements ON mission_requirements.id = mission_progresses.mission_requirement_id JOIN missions ON missions.id = mission_requirements.mission_id WHERE (mission_requirements.id = ? AND mission_progresses.user_id = ?) AND missions.status NOT IN (?, ?) AND (`mission_progresses`.`user_id` = ? AND `mission_progresses`.`mission_requirement_id` = ?) AND `mission_progresses`.`deleted_at` IS NULL ORDER BY `mission_progresses`.`id` LIMIT ?")).
					WithArgs(1, 1, domain.MissionUnavailable, domain.MissionCanceled, 1, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "mission_requirement_id", "progress", "completed", "created_at", "updated_at"}).
						AddRow(1, 1, 1, 0, false, fixedTime, fixedTime))
			},
			expectErr: false,
			expectProgress: &domain.MissionProgress{
				ID:                   1,
				Progress:             0,
				Completed:            false,
				UserID:               1,
				MissionRequirementID: 1,
			},
		},
		"create new progress": {
			userID:               1,
			missionRequirementID: 2,
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT `mission_progresses`.`id`,`mission_progresses`.`created_at`,`mission_progresses`.`updated_at`,`mission_progresses`.`deleted_at`,`mission_progresses`.`progress`,`mission_progresses`.`completed`,`mission_progresses`.`user_id`,`mission_progresses`.`mission_requirement_id` FROM `mission_progresses` JOIN mission_requirements ON mission_requirements.id = mission_progresses.mission_requirement_id JOIN missions ON missions.id = mission_requirements.mission_id WHERE (mission_requirements.id = ? AND mission_progresses.user_id = ?) AND missions.status NOT IN (?, ?) AND (`mission_progresses`.`user_id` = ? AND `mission_progresses`.`mission_requirement_id` = ?) AND `mission_progresses`.`deleted_at` IS NULL ORDER BY `mission_progresses`.`id` LIMIT ?")).
					WithArgs(2, 1, domain.TitleUnavailable, domain.TitleCanceled, 1, 2, 1).
					WillReturnRows(sqlmock.NewRows([]string{}))

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `mission_progresses` (`created_at`,`updated_at`,`deleted_at`,`progress`,`completed`,`user_id`,`mission_requirement_id`) VALUES (?,?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 0, false, 1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
			expectProgress: &domain.MissionProgress{
				ID:                   1,
				Progress:             0,
				Completed:            false,
				UserID:               1,
				MissionRequirementID: 2,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mock()

			progress, err := r.GetOrCreateMissionProgress(tc.userID, tc.missionRequirementID)

			if (err != nil) != tc.expectErr {
				t.Errorf("expected error: %v, got: %v", tc.expectErr, err)
			} else {
				assert.Equal(t, progress.ID, tc.expectProgress.ID)
				assert.Equal(t, progress.Progress, tc.expectProgress.Progress)
				assert.Equal(t, progress.Completed, tc.expectProgress.Completed)
				assert.Equal(t, progress.UserID, tc.expectProgress.UserID)
				assert.Equal(t, progress.MissionRequirementID, tc.expectProgress.MissionRequirementID)

				assert.NotZero(t, progress.CreatedAt)
				assert.NotZero(t, progress.UpdatedAt)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestTaskRepositoryMySQL_AwardTitleToUser(t *testing.T) {
	gormDB, mock := tests.Setup(t)
	r := db.NewTaskRepositoryMySQL(gormDB)

	userID := uint(1)
	titleID := uint(1)

	testCases := map[string]struct {
		mock      func()
		expectErr bool
	}{
		"successful award new title": {
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_titles`")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), false, userID, titleID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `title_requirements` WHERE title_id = ? AND `title_requirements`.`deleted_at` IS NULL")).
					WithArgs(titleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title_id", "goal"}).
						AddRow(1, titleID, 100))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `title_progresses` WHERE (user_id = ? AND title_requirement_id = ?) AND `title_progresses`.`deleted_at` IS NULL ORDER BY `title_progresses`.`id` LIMIT ?")).
					WithArgs(userID, 1, 1).
					WillReturnError(gorm.ErrRecordNotFound)

				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `title_progresses`")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 100, true, userID, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectErr: false,
		},
		"error on creating user title": {
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_titles`")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), false, userID, titleID).
					WillReturnError(errors.New("insert error"))

				mock.ExpectRollback()
			},
			expectErr: true,
		},
		"error on finding requirements": {
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_titles`")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), false, userID, titleID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `title_requirements` WHERE title_id = ? AND `title_requirements`.`deleted_at` IS NULL")).
					WithArgs(titleID).
					WillReturnError(errors.New("query error"))

				mock.ExpectRollback()
			},
			expectErr: true,
		},
		"error on creating progress": {
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_titles`")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), false, userID, titleID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `title_requirements` WHERE title_id = ? AND `title_requirements`.`deleted_at` IS NULL")).
					WithArgs(titleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title_id", "goal"}).
						AddRow(1, titleID, 100))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `title_progresses` WHERE (user_id = ? AND title_requirement_id = ?) AND `title_progresses`.`deleted_at` IS NULL ORDER BY `title_progresses`.`id` LIMIT ?")).
					WithArgs(userID, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "title_requirement_id", "progress", "completed"}).
						AddRow(userID, 1, 50, false))

				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `title_progresses`")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 100, true, userID, 1).
					WillReturnError(errors.New("insert progress error"))

				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mock()

			err := r.AwardTitleToUser(userID, titleID)

			if (err != nil) != tc.expectErr {
				t.Errorf("expected error: %v, got: %v", tc.expectErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
