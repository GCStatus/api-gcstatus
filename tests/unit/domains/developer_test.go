package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
	"gcstatus/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateDeveloper(t *testing.T) {
	testCases := map[string]struct {
		developer    domain.Developer
		mockBehavior func(mock sqlmock.Sqlmock, developer domain.Developer)
		expectError  bool
	}{
		"Success": {
			developer: domain.Developer{
				Name:   "Game Science",
				Acting: false,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, developer domain.Developer) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `developers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						developer.Name,
						developer.Acting,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			developer: domain.Developer{
				Name:   "Game Science",
				Acting: false,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, developer domain.Developer) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `developers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						developer.Name,
						developer.Acting,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := tests.Setup(t)

			tc.mockBehavior(mock, tc.developer)

			err := db.Create(&tc.developer).Error

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

func TestUpdateDeveloper(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		developer    domain.Developer
		mockBehavior func(mock sqlmock.Sqlmock, developer domain.Developer)
		expectError  bool
	}{
		"Success": {
			developer: domain.Developer{
				ID:        1,
				Name:      "Game Science",
				Acting:    true,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, developer domain.Developer) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `developers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						developer.Name,
						developer.Acting,
						developer.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			developer: domain.Developer{
				ID:        1,
				Name:      "Game Science",
				Acting:    false,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, developer domain.Developer) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `developers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						developer.Name,
						developer.Acting,
						developer.ID,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := tests.Setup(t)

			tc.mockBehavior(mock, tc.developer)

			err := db.Save(&tc.developer).Error

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

func TestSoftDeleteDeveloper(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		developerID  uint
		mockBehavior func(mock sqlmock.Sqlmock, developerID uint)
		wantErr      bool
	}{
		"Can soft delete a Developer": {
			developerID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, developerID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `developers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), developerID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			developerID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, developerID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `developers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Developer"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.developerID)

			err := db.Delete(&domain.Developer{}, tc.developerID).Error

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

func TestValidateDeveloper(t *testing.T) {
	testCases := map[string]struct {
		developer domain.Developer
	}{
		"Can empty validations errors": {
			developer: domain.Developer{
				Name:   "Game Science",
				Acting: true,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.developer.ValidateDeveloper()
			assert.NoError(t, err)
		})
	}
}

func TestCreateDeveloperWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		developer domain.Developer
		wantErr   string
	}{
		"Missing required fields": {
			developer: domain.Developer{},
			wantErr: `
				Name is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.developer.ValidateDeveloper()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
