package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	testutils "gcstatus/tests/utils"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func CreateTitleProgressTest(t *testing.T) {
	testCases := map[string]struct {
		titleProgress domain.TitleProgress
		mockBehavior  func(mock sqlmock.Sqlmock, titleProgress domain.TitleProgress)
		expectErr     bool
	}{
		"Successfully created": {
			titleProgress: domain.TitleProgress{
				Progress:  5,
				Completed: false,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, titleProgress domain.TitleProgress) {
				mock.ExpectBegin()

				mock.ExpectExec("^INSERT INTO `title_progresses`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						titleProgress.Progress,
						titleProgress.Completed,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		"Failure - Insert Error": {
			titleProgress: domain.TitleProgress{
				Progress:  5,
				Completed: false,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, titleProgress domain.TitleProgress) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `title_progresses`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						titleProgress.Progress,
						titleProgress.Completed,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := testutils.Setup(t)

			tc.mockBehavior(mock, tc.titleProgress)

			err := db.Create(&tc.titleProgress).Error

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSoftDeleteTitleProgress(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		titleProgressID uint
		mockBehavior    func(mock sqlmock.Sqlmock, titleProgressID uint)
		wantErr         bool
	}{
		"Can soft delete a title progress": {
			titleProgressID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, titleProgressID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `title_progresses` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), titleProgressID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			titleProgressID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, titleProgressID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `title_progresses` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete title requirement"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.titleProgressID)

			err := db.Delete(&domain.TitleProgress{}, tc.titleProgressID).Error

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateTitleProgress(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		titleProgress domain.TitleProgress
		mockBehavior  func(mock sqlmock.Sqlmock, titleProgress domain.TitleProgress)
		expectError   bool
	}{
		"Success": {
			titleProgress: domain.TitleProgress{
				ID:                 1,
				Progress:           5,
				Completed:          false,
				CreatedAt:          fixedTime,
				UpdatedAt:          fixedTime,
				UserID:             1,
				TitleRequirementID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, titleProgress domain.TitleProgress) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `title_progresses`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						titleProgress.Progress,
						titleProgress.Completed,
						titleProgress.UserID,
						titleProgress.TitleRequirementID,
						titleProgress.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			titleProgress: domain.TitleProgress{
				ID:                 1,
				Progress:           5,
				Completed:          false,
				CreatedAt:          fixedTime,
				UpdatedAt:          fixedTime,
				UserID:             1,
				TitleRequirementID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, titleProgress domain.TitleProgress) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `title_progresses`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						titleProgress.Progress,
						titleProgress.Completed,
						titleProgress.UserID,
						titleProgress.TitleRequirementID,
						titleProgress.ID,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := testutils.Setup(t)

			tc.mockBehavior(mock, tc.titleProgress)

			err := db.Save(&tc.titleProgress).Error

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidateTitleProgressValidData(t *testing.T) {
	testCases := map[string]struct {
		titleProgress domain.TitleProgress
	}{
		"Can empty validations errors": {
			titleProgress: domain.TitleProgress{
				Progress:  5,
				Completed: false,
				User: domain.User{
					Name:       "Name",
					Email:      "test@example.com",
					Nickname:   "test1",
					Experience: 100,
					Birthdate:  time.Now(),
					Password:   "fakepass123",
					Profile: domain.Profile{
						Share: true,
					},
					Level: domain.Level{
						Level:      1,
						Coins:      100,
						Experience: 100,
					},
					Wallet: domain.Wallet{
						Amount: 100,
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.titleProgress.ValidateTitleProgress()

			assert.NoError(t, err)
		})
	}
}

func TestCreateTitleProgressWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		titleProgress domain.TitleProgress
		wantErr       string
	}{
		"Missing required fields": {
			titleProgress: domain.TitleProgress{},
			wantErr: `
				Name is a required field,
				Email is a required field,
				Nickname is a required field,
				Birthdate is a required field,
				Password is a required field,
				Share is a required field,
				Level is a required field,
				Experience is a required field,
				Coins is a required field,
				Amount is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.titleProgress.ValidateTitleProgress()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
