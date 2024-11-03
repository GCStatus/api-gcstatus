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

func TestCreateCracker(t *testing.T) {
	testCases := map[string]struct {
		cracker      domain.Cracker
		mockBehavior func(mock sqlmock.Sqlmock, cracker domain.Cracker)
		expectError  bool
	}{
		"Success": {
			cracker: domain.Cracker{
				Name:   "GOLDBERG",
				Slug:   "goldberg",
				Acting: false,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, cracker domain.Cracker) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `crackers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						cracker.Name,
						cracker.Slug,
						cracker.Acting,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			cracker: domain.Cracker{
				Name:   "GOLDBERG",
				Slug:   "goldberg",
				Acting: false,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, cracker domain.Cracker) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `crackers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						cracker.Name,
						cracker.Slug,
						cracker.Acting,
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

			tc.mockBehavior(mock, tc.cracker)

			err := db.Create(&tc.cracker).Error

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

func TestUpdateCracker(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		cracker      domain.Cracker
		mockBehavior func(mock sqlmock.Sqlmock, cracker domain.Cracker)
		expectError  bool
	}{
		"Success": {
			cracker: domain.Cracker{
				ID:        1,
				Name:      "GOLDBERG",
				Slug:      "goldberg",
				Acting:    true,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, cracker domain.Cracker) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `crackers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						cracker.Name,
						cracker.Slug,
						cracker.Acting,
						cracker.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			cracker: domain.Cracker{
				ID:        1,
				Name:      "GOLDBERG",
				Slug:      "goldberg",
				Acting:    true,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, cracker domain.Cracker) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `crackers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						cracker.Name,
						cracker.Slug,
						cracker.Acting,
						cracker.ID,
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

			tc.mockBehavior(mock, tc.cracker)

			err := db.Save(&tc.cracker).Error

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

func TestSoftDeleteCracker(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		crackerID    uint
		mockBehavior func(mock sqlmock.Sqlmock, crackerID uint)
		wantErr      bool
	}{
		"Can soft delete a Cracker": {
			crackerID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, crackerID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `crackers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), crackerID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			crackerID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, crackerID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `crackers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Cracker"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.crackerID)

			err := db.Delete(&domain.Cracker{}, tc.crackerID).Error

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

func TestValidateCracker(t *testing.T) {
	testCases := map[string]struct {
		cracker domain.Cracker
	}{
		"Can empty validations errors": {
			cracker: domain.Cracker{
				Name:   "GOLDBERG",
				Slug:   "goldberg",
				Acting: true,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.cracker.ValidateCracker()
			assert.NoError(t, err)
		})
	}
}

func TestCreateCrackerWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		cracker domain.Cracker
		wantErr string
	}{
		"Missing required fields": {
			cracker: domain.Cracker{},
			wantErr: `
				Name is a required field,
				Slug is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.cracker.ValidateCracker()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
