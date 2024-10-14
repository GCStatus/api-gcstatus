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

func CreateTitleRequirementTest(t *testing.T) {
	testCases := map[string]struct {
		titleRequirement domain.TitleRequirement
		mockBehavior     func(mock sqlmock.Sqlmock, title domain.TitleRequirement)
		expectErr        bool
	}{
		"Successfully created": {
			titleRequirement: domain.TitleRequirement{
				Task:        "Do something",
				Key:         "do_something",
				Goal:        10,
				Description: "Title 1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, titleRequirement domain.TitleRequirement) {
				mock.ExpectBegin()

				mock.ExpectExec("^INSERT INTO `title_requirements`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						titleRequirement.Task,
						titleRequirement.Description,
						titleRequirement.Key,
						titleRequirement.Goal,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		"Failure - Insert Error": {
			titleRequirement: domain.TitleRequirement{
				Task:        "Do something",
				Key:         "do_something",
				Goal:        10,
				Description: "Title 1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, titleRequirement domain.TitleRequirement) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `title_requirements`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						titleRequirement.Task,
						titleRequirement.Description,
						titleRequirement.Key,
						titleRequirement.Goal,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := tests.Setup(t)

			tc.mockBehavior(mock, tc.titleRequirement)

			err := db.Create(&tc.titleRequirement).Error

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

func TestSoftDeleteTitleRequirement(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		titleRequirementID uint
		mockBehavior       func(mock sqlmock.Sqlmock, titleRequirementID uint)
		wantErr            bool
	}{
		"Can soft delete a title requirement": {
			titleRequirementID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, titleRequirementID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `title_requirements` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), titleRequirementID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			titleRequirementID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, titleRequirementID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `title_requirements` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete title requirement"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.titleRequirementID)

			err := db.Delete(&domain.TitleRequirement{}, tc.titleRequirementID).Error

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

func TestUpdateTitleRequirement(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		titleRequirement domain.TitleRequirement
		mockBehavior     func(mock sqlmock.Sqlmock, title domain.TitleRequirement)
		expectError      bool
	}{
		"Success": {
			titleRequirement: domain.TitleRequirement{
				ID:          1,
				Task:        "Do something",
				Key:         "do_something",
				Description: "Title 1",
				Goal:        10,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				TitleID:     1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, titleRequirement domain.TitleRequirement) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `title_requirements`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						titleRequirement.Task,
						titleRequirement.Key,
						titleRequirement.Goal,
						titleRequirement.Description,
						titleRequirement.TitleID,
						titleRequirement.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			titleRequirement: domain.TitleRequirement{
				ID:          1,
				Task:        "Do something",
				Key:         "do_something",
				Description: "Title 1",
				Goal:        10,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
				TitleID:     1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, titleRequirement domain.TitleRequirement) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `title_requirements`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						titleRequirement.Task,
						titleRequirement.Key,
						titleRequirement.Goal,
						titleRequirement.Description,
						titleRequirement.TitleID,
						titleRequirement.ID,
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

			tc.mockBehavior(mock, tc.titleRequirement)

			err := db.Save(&tc.titleRequirement).Error

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

func TestValidateTitleRequirementValidData(t *testing.T) {
	testCases := map[string]struct {
		titleRequirement domain.TitleRequirement
	}{
		"Can empty validations errors": {
			titleRequirement: domain.TitleRequirement{
				ID:          1,
				Task:        "Do something",
				Key:         "do_something",
				Description: "Title 1",
				Goal:        10,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				TitleID:     1,
				TitleProgress: domain.TitleProgress{
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
				Title: domain.Title{
					Title:       "Title 1",
					Description: "Title 1",
					Purchasable: false,
					Status:      "available",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.titleRequirement.ValidateTitleRequirement()

			assert.NoError(t, err)
		})
	}
}

func TestCreateTitleRequirementWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		titleRequirement domain.TitleRequirement
		wantErr          string
	}{
		"Missing required fields": {
			titleRequirement: domain.TitleRequirement{},
			wantErr: `
				Task is a required field,
				Key is a required field,
				Goal is a required field,
				Description is a required field,
				Name is a required field,
				Email is a required field,
				Nickname is a required field,
				Birthdate is a required field,
				Password is a required field,
				Share is a required field,
				Level is a required field,
				Experience is a required field,
				Coins is a required field,
				Amount is a required field,
				Title is a required field,
				Description is a required field,
				Status is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.titleRequirement.ValidateTitleRequirement()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
