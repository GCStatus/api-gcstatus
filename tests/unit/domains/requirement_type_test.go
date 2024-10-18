package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateRequirementType(t *testing.T) {
	testCases := map[string]struct {
		requirementType domain.RequirementType
		mockBehavior    func(mock sqlmock.Sqlmock, requirementType domain.RequirementType)
		expectError     bool
	}{
		"Success": {
			requirementType: domain.RequirementType{
				Potential: domain.MinimumRequirementType,
				OS:        domain.WindowsOSRequirement,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, requirementType domain.RequirementType) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `requirement_types`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						requirementType.Potential,
						requirementType.OS,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			requirementType: domain.RequirementType{
				Potential: domain.MinimumRequirementType,
				OS:        domain.WindowsOSRequirement,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, requirementType domain.RequirementType) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `requirement_types`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						requirementType.Potential,
						requirementType.OS,
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

			tc.mockBehavior(mock, tc.requirementType)

			err := db.Create(&tc.requirementType).Error

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

func TestUpdateRequirementType(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		requirementType domain.RequirementType
		mockBehavior    func(mock sqlmock.Sqlmock, requirementType domain.RequirementType)
		expectError     bool
	}{
		"Success": {
			requirementType: domain.RequirementType{
				ID:        1,
				Potential: domain.MinimumRequirementType,
				OS:        domain.WindowsOSRequirement,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, requirementType domain.RequirementType) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `requirement_types`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						requirementType.Potential,
						requirementType.OS,
						requirementType.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			requirementType: domain.RequirementType{
				ID:        1,
				Potential: domain.MinimumRequirementType,
				OS:        domain.WindowsOSRequirement,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, requirementType domain.RequirementType) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `requirement_types`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						requirementType.Potential,
						requirementType.OS,
						requirementType.ID,
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

			tc.mockBehavior(mock, tc.requirementType)

			err := db.Save(&tc.requirementType).Error

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

func TestSoftDeleteRequirementType(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		requirementTypeID uint
		mockBehavior      func(mock sqlmock.Sqlmock, requirementTypeID uint)
		wantErr           bool
	}{
		"Can soft delete a RequirementType": {
			requirementTypeID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, requirementTypeID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `requirement_types` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), requirementTypeID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			requirementTypeID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, requirementTypeID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `requirement_types` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete RequirementType"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.requirementTypeID)

			err := db.Delete(&domain.RequirementType{}, tc.requirementTypeID).Error

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

func TestValidateRequirementTypeLanguageValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		requirementType domain.RequirementType
	}{
		"Can empty validations errors": {
			requirementType: domain.RequirementType{
				Potential: domain.MinimumRequirementType,
				OS:        domain.WindowsOSRequirement,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.requirementType.ValidateRequirementType()
			assert.NoError(t, err)
		})
	}
}

func TestCreateRequirementTypeWithMissingFields(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		requirementType domain.RequirementType
		wantErr         string
	}{
		"Missing required fields": {
			requirementType: domain.RequirementType{},
			wantErr:         "Potential is a required field, OS is a required field",
		},
		"Invalid potential": {
			requirementType: domain.RequirementType{
				Potential: "invalid",
				OS:        domain.WindowsOSRequirement,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			wantErr: "Potential must be one of 'minimum', 'recommended', or 'maximum'",
		},
		"invalid os": {
			requirementType: domain.RequirementType{
				Potential: domain.MinimumRequirementType,
				OS:        "invalid",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			wantErr: "OS must be one of 'windows', 'mac', or 'linux'",
		},
		"both invalid": {
			requirementType: domain.RequirementType{
				Potential: "invalid",
				OS:        "invalid",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			wantErr: "Potential must be one of 'minimum', 'recommended', or 'maximum', OS must be one of 'windows', 'mac', or 'linux'",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.requirementType.ValidateRequirementType()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
