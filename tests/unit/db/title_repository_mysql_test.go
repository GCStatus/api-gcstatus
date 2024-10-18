package tests

import (
	"errors"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTitleRepositoryMySQL_GetAllForUser(t *testing.T) {
	gormDB, mock := testutils.Setup(t)
	repo := db.NewTitleRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		userID         uint
		mockSetup      func()
		expectedError  error
		expectedTitles []domain.Title
	}{
		"success - titles found": {
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `titles` WHERE status NOT IN (?, ?) AND `titles`.`deleted_at` IS NULL")).
					WithArgs(domain.TitleUnavailable, domain.TitleCanceled).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "status"}).
						AddRow(1, "Title 1", "Available").
						AddRow(2, "Title 2", "Available"))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `title_requirements` WHERE `title_requirements`.`title_id` IN (?,?) AND `title_requirements`.`deleted_at` IS NULL")).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "task", "key", "goal", "description", "title_id"}).
						AddRow(1, "Task 1", "Key 1", 100, "Description 1", 1).
						AddRow(2, "Task 2", "Key 2", 200, "Description 2", 2))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `title_progresses` WHERE `title_progresses`.`title_requirement_id` IN (?,?) AND user_id = ? AND `title_progresses`.`deleted_at` IS NULL")).
					WithArgs(1, 2, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "progress", "title_requirement_id", "user_id"}).
						AddRow(1, 50, 1, 1).
						AddRow(2, 75, 2, 1))
			},
			expectedTitles: []domain.Title{
				{
					ID:     1,
					Title:  "Title 1",
					Status: "Available",
					TitleRequirements: []domain.TitleRequirement{
						{
							ID:            1,
							Task:          "Task 1",
							Key:           "Key 1",
							Goal:          100,
							Description:   "Description 1",
							TitleID:       1,
							TitleProgress: domain.TitleProgress{ID: 1, Progress: 50, UserID: 1, TitleRequirementID: 1},
						},
					},
				},
				{
					ID:     2,
					Title:  "Title 2",
					Status: "Available",
					TitleRequirements: []domain.TitleRequirement{
						{
							ID:            2,
							Task:          "Task 2",
							Key:           "Key 2",
							Goal:          200,
							Description:   "Description 2",
							TitleID:       2,
							TitleProgress: domain.TitleProgress{ID: 2, Progress: 75, UserID: 1, TitleRequirementID: 2},
						},
					},
				},
			},
			expectedError: nil,
		},
		"no titles found": {
			userID: 2,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `titles` WHERE status NOT IN (?, ?) AND `titles`.`deleted_at` IS NULL")).
					WithArgs(domain.TitleUnavailable, domain.TitleCanceled).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectedTitles: []domain.Title{},
			expectedError:  nil,
		},
		"error - db failure": {
			userID: 3,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `titles` WHERE status NOT IN (?, ?) AND `titles`.`deleted_at` IS NULL")).
					WithArgs(domain.TitleUnavailable, domain.TitleCanceled).
					WillReturnError(errors.New("db error"))
			},
			expectedTitles: nil,
			expectedError:  errors.New("db error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()

			titles, err := repo.GetAllForUser(tc.userID)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedTitles, titles)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTitleRepositoryMySQL_FindById(t *testing.T) {
	gormDB, mock := testutils.Setup(t)
	repo := db.NewTitleRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		titleID       uint
		mockSetup     func()
		expectedError error
		expectedTitle domain.Title
	}{
		"success - title found": {
			titleID: 1,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `titles` WHERE id = ? AND `titles`.`deleted_at` IS NULL ORDER BY `titles`.`id` LIMIT ?")).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "status"}).
						AddRow(1, "Title 1", "Available"))
			},
			expectedTitle: domain.Title{
				ID:     1,
				Title:  "Title 1",
				Status: "Available",
			},
			expectedError: nil,
		},
		"error - title not found": {
			titleID: 2,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `titles` WHERE id = ? AND `titles`.`deleted_at` IS NULL ORDER BY `titles`.`id` LIMIT ?")).
					WithArgs(2, 1).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectedTitle: domain.Title{},
			expectedError: gorm.ErrRecordNotFound,
		},
		"error - db failure": {
			titleID: 3,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `titles` WHERE id = ? AND `titles`.`deleted_at` IS NULL ORDER BY `titles`.`id` LIMIT ?")).
					WithArgs(3, 1).
					WillReturnError(errors.New("db error"))
			},
			expectedTitle: domain.Title{},
			expectedError: errors.New("db error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()

			title, err := repo.FindById(tc.titleID)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedTitle, title)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTitleRepositoryMySQL_ToggleEnableTitle(t *testing.T) {
	gormDB, mock := testutils.Setup(t)
	repo := db.NewTitleRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		userID        uint
		titleID       uint
		mockSetup     func()
		expectedError error
	}{
		"success - toggle enable title": {
			userID:  1,
			titleID: 1,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_titles` WHERE (user_id = ? AND title_id = ?) AND `user_titles`.`deleted_at` IS NULL ORDER BY `user_titles`.`id` LIMIT ?")).
					WithArgs(1, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title_id", "enabled"}).
						AddRow(1, 1, 1, false))

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_titles` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`enabled`=?,`user_id`=?,`title_id`=? WHERE `user_titles`.`deleted_at` IS NULL AND `id` = ?")).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						true,
						1,
						1,
						1,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_titles` SET `enabled`=?,`updated_at`=? WHERE (user_id = ? AND title_id != ?) AND `user_titles`.`deleted_at` IS NULL")).
					WithArgs(
						false,
						sqlmock.AnyArg(),
						1,
						1,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		"error - title not found": {
			userID:  2,
			titleID: 2,
			mockSetup: func() {

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_titles` WHERE (user_id = ? AND title_id = ?) AND `user_titles`.`deleted_at` IS NULL ORDER BY `user_titles`.`id` LIMIT ?")).
					WithArgs(2, 2, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
		},
		"error - updating other titles": {
			userID:  4,
			titleID: 4,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_titles` WHERE (user_id = ? AND title_id = ?) AND `user_titles`.`deleted_at` IS NULL ORDER BY `user_titles`.`id` LIMIT ?")).
					WithArgs(4, 4, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title_id", "enabled"}).
						AddRow(1, 4, 4, false))

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_titles` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`enabled`=?,`user_id`=?,`title_id`=? WHERE `user_titles`.`deleted_at` IS NULL AND `id` = ?")).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						true,
						4,
						4,
						1,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `user_titles` SET `enabled`=?,`updated_at`=? WHERE (user_id = ? AND title_id != ?) AND `user_titles`.`deleted_at` IS NULL")).
					WithArgs(
						false,
						sqlmock.AnyArg(),
						4,
						4,
					).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("db error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			tc.mockSetup()

			err := repo.ToggleEnableTitle(tc.userID, tc.titleID)

			assert.Equal(t, tc.expectedError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
