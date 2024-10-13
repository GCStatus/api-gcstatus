package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"gcstatus/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateCrack(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		crack        domain.Crack
		mockBehavior func(mock sqlmock.Sqlmock, crack domain.Crack)
		expectError  bool
	}{
		"Success": {
			crack: domain.Crack{
				CrackedAt:    utils.TimePtr(fixedTime),
				Status:       domain.UncrackedStatus,
				CrackerID:    1,
				ProtectionID: 1,
				GameID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, crack domain.Crack) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `cracks`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						crack.Status,
						crack.CrackedAt,
						crack.CrackerID,
						crack.ProtectionID,
						crack.GameID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			crack: domain.Crack{
				Status:       domain.UncrackedStatus,
				CrackerID:    1,
				ProtectionID: 1,
				GameID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, crack domain.Crack) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `cracks`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						crack.Status,
						crack.CrackedAt,
						crack.CrackerID,
						crack.ProtectionID,
						crack.GameID,
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

			tc.mockBehavior(mock, tc.crack)

			err := db.Create(&tc.crack).Error

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

func TestUpdateCrack(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		crack        domain.Crack
		mockBehavior func(mock sqlmock.Sqlmock, crack domain.Crack)
		expectError  bool
	}{
		"Success": {
			crack: domain.Crack{
				ID:           1,
				Status:       domain.UncrackedStatus,
				CrackerID:    1,
				ProtectionID: 1,
				GameID:       1,
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, crack domain.Crack) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `cracks`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						crack.Status,
						crack.CrackedAt,
						crack.CrackerID,
						crack.ProtectionID,
						crack.GameID,
						crack.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			crack: domain.Crack{
				ID:           1,
				Status:       domain.UncrackedStatus,
				CrackerID:    1,
				ProtectionID: 1,
				GameID:       1,
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, crack domain.Crack) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `cracks`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						crack.Status,
						crack.CrackedAt,
						crack.CrackerID,
						crack.ProtectionID,
						crack.GameID,
						crack.ID,
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

			tc.mockBehavior(mock, tc.crack)

			err := db.Save(&tc.crack).Error

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

func TestSoftDeleteCrack(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		crackID      uint
		mockBehavior func(mock sqlmock.Sqlmock, crackID uint)
		wantErr      bool
	}{
		"Can soft delete a Crack": {
			crackID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, crackID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `cracks` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), crackID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			crackID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, crackID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `cracks` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Crack"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.crackID)

			err := db.Delete(&domain.Crack{}, tc.crackID).Error

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

func TestValidateCrack(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		crack domain.Crack
	}{
		"Can empty validations errors": {
			crack: domain.Crack{
				Status: domain.UncrackedStatus,
				Cracker: domain.Cracker{
					Name:   "GOLDBERG",
					Acting: true,
				},
				Protection: domain.Protection{
					Name: "Denuvo",
				},
				Game: domain.Game{
					Slug:             "valid",
					Age:              18,
					Title:            "Game Test",
					Condition:        domain.CommomCondition,
					Cover:            "https://placehold.co/600x400/EEE/31343C",
					About:            "About game",
					Description:      "Description",
					ShortDescription: "Short description",
					Free:             false,
					ReleaseDate:      fixedTime,
					CreatedAt:        fixedTime,
					UpdatedAt:        fixedTime,
					Views: []domain.Viewable{
						{
							UserID:       10,
							ViewableID:   1,
							ViewableType: "games",
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.crack.ValidateCrack()
			assert.NoError(t, err)
		})
	}
}

func TestCreateCrackWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		crack   domain.Crack
		wantErr string
	}{
		"Missing required fields": {
			crack: domain.Crack{},
			wantErr: `
				Status is a required field,
				Name is a required field,
				Name is a required field,
				Age is a required field,
				Slug is a required field,
				Title is a required field,
				Condition is a required field,
				Cover is a required field,
				About is a required field,
				Description is a required field,
				ShortDescription is a required field,
				ReleaseDate is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.crack.ValidateCrack()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
