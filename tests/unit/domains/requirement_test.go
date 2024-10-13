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

func TestCreateRequirement(t *testing.T) {
	testCases := map[string]struct {
		requirement  domain.Requirement
		mockBehavior func(mock sqlmock.Sqlmock, requirement domain.Requirement)
		expectError  bool
	}{
		"Success": {
			requirement: domain.Requirement{
				OS:                "Windows 11 64 bits",
				DX:                "DirectX 12",
				CPU:               "Ryzen 5 3600",
				RAM:               "16GB",
				GPU:               "RTX 3090 TI",
				ROM:               "90GB",
				OBS:               utils.StringPtr("Some observation"),
				Network:           "Non required",
				RequirementTypeID: 1,
				GameID:            1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, requirement domain.Requirement) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `requirements`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						requirement.OS,
						requirement.DX,
						requirement.CPU,
						requirement.RAM,
						requirement.GPU,
						requirement.ROM,
						requirement.OBS,
						requirement.Network,
						requirement.RequirementTypeID,
						requirement.GameID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			requirement: domain.Requirement{
				OS:                "Windows 11 64 bits",
				DX:                "DirectX 12",
				CPU:               "Ryzen 5 3600",
				RAM:               "16GB",
				GPU:               "RTX 3090 TI",
				ROM:               "90GB",
				OBS:               utils.StringPtr("Some observation"),
				Network:           "Non required",
				RequirementTypeID: 1,
				GameID:            1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, requirement domain.Requirement) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `requirements`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						requirement.OS,
						requirement.DX,
						requirement.CPU,
						requirement.RAM,
						requirement.GPU,
						requirement.ROM,
						requirement.OBS,
						requirement.Network,
						requirement.RequirementTypeID,
						requirement.GameID,
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

			tc.mockBehavior(mock, tc.requirement)

			err := db.Create(&tc.requirement).Error

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

func TestUpdateRequirement(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		requirement  domain.Requirement
		mockBehavior func(mock sqlmock.Sqlmock, requirement domain.Requirement)
		expectError  bool
	}{
		"Success": {
			requirement: domain.Requirement{
				ID:                1,
				OS:                "Windows 11 64 bits",
				DX:                "DirectX 12",
				CPU:               "Ryzen 5 3600",
				RAM:               "16GB",
				GPU:               "RTX 3090 TI",
				ROM:               "90GB",
				OBS:               utils.StringPtr("Some observation"),
				Network:           "Non required",
				RequirementTypeID: 1,
				GameID:            1,
				CreatedAt:         fixedTime,
				UpdatedAt:         fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, requirement domain.Requirement) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `requirements`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						requirement.OS,
						requirement.DX,
						requirement.CPU,
						requirement.RAM,
						requirement.GPU,
						requirement.ROM,
						requirement.OBS,
						requirement.Network,
						requirement.RequirementTypeID,
						requirement.GameID,
						requirement.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			requirement: domain.Requirement{
				ID:                1,
				OS:                "Windows 11 64 bits",
				DX:                "DirectX 12",
				CPU:               "Ryzen 5 3600",
				RAM:               "16GB",
				GPU:               "RTX 3090 TI",
				ROM:               "90GB",
				OBS:               utils.StringPtr("Some observation"),
				Network:           "Non required",
				RequirementTypeID: 1,
				GameID:            1,
				CreatedAt:         fixedTime,
				UpdatedAt:         fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, requirement domain.Requirement) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `requirements`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						requirement.OS,
						requirement.DX,
						requirement.CPU,
						requirement.RAM,
						requirement.GPU,
						requirement.ROM,
						requirement.OBS,
						requirement.Network,
						requirement.RequirementTypeID,
						requirement.GameID,
						requirement.ID,
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

			tc.mockBehavior(mock, tc.requirement)

			err := db.Save(&tc.requirement).Error

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

func TestSoftDeleteRequirement(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		requirementID uint
		mockBehavior  func(mock sqlmock.Sqlmock, requirementID uint)
		wantErr       bool
	}{
		"Can soft delete a Requirement": {
			requirementID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, requirementID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `requirements` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), requirementID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			requirementID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, requirementID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `requirements` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Requirement"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.requirementID)

			err := db.Delete(&domain.Requirement{}, tc.requirementID).Error

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

func TestValidateRequirement(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		requirement domain.Requirement
	}{
		"Can empty validations errors": {
			requirement: domain.Requirement{
				OS:      "Windows 11 64 bits",
				DX:      "DirectX 12",
				CPU:     "Ryzen 5 3600",
				RAM:     "16GB",
				GPU:     "RTX 3090 TI",
				ROM:     "90GB",
				OBS:     utils.StringPtr("Some observation"),
				Network: "Non required",
				RequirementType: domain.RequirementType{
					Potential: domain.MinimumRequirementType,
					OS:        domain.WindowsOSRequirement,
				},
				Game: domain.Game{
					ID:               1,
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
			err := tc.requirement.ValidateRequirement()
			assert.NoError(t, err)
		})
	}
}

func TestCreateRequirementWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		requirement domain.Requirement
		wantErr     string
	}{
		"Missing required fields": {
			requirement: domain.Requirement{},
			wantErr: `
				OS is a required field,
				DX is a required field,
				CPU is a required field,
				RAM is a required field,
				GPU is a required field,
				ROM is a required field,
				Network is a required field,
				Potential is a required field,
				OS is a required field,
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
			err := tc.requirement.ValidateRequirement()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
