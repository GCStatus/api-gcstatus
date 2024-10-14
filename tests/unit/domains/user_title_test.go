package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"gcstatus/tests"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func CreateUserTitleTest(t *testing.T) {
	testCases := map[string]struct {
		userTitle    domain.UserTitle
		mockBehavior func(mock sqlmock.Sqlmock, userTitle domain.UserTitle)
		expectErr    bool
	}{
		"Successfully created": {
			userTitle: domain.UserTitle{
				Enabled: false,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, userTitle domain.UserTitle) {
				mock.ExpectBegin()

				mock.ExpectExec("^INSERT INTO `user_titles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						userTitle.Enabled,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		"Failure - Insert Error": {
			userTitle: domain.UserTitle{
				Enabled: false,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, userTitle domain.UserTitle) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `user_titles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						userTitle.Enabled,
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

			tc.mockBehavior(mock, tc.userTitle)

			err := db.Create(&tc.userTitle).Error

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

func TestSoftDeleteUserTitle(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		userTitleID  uint
		mockBehavior func(mock sqlmock.Sqlmock, userTitleID uint)
		wantErr      bool
	}{
		"Can soft delete a title": {
			userTitleID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, userTitleID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user_titles` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), userTitleID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			userTitleID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, userTitleID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user_titles` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete title"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.userTitleID)

			err := db.Delete(&domain.UserTitle{}, tc.userTitleID).Error

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

func TestUpdateUserTitle(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		userTitle    domain.UserTitle
		mockBehavior func(mock sqlmock.Sqlmock, userTitle domain.UserTitle)
		expectError  bool
	}{
		"Success": {
			userTitle: domain.UserTitle{
				ID:        1,
				Enabled:   true,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, userTitle domain.UserTitle) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user_titles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						userTitle.Enabled,
						userTitle.UserID,
						userTitle.TitleID,
						userTitle.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			userTitle: domain.UserTitle{
				ID:        1,
				Enabled:   false,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, userTitle domain.UserTitle) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `user_titles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						userTitle.Enabled,
						userTitle.UserID,
						userTitle.TitleID,
						userTitle.ID,
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

			tc.mockBehavior(mock, tc.userTitle)

			err := db.Save(&tc.userTitle).Error

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

func TestValidateUserTitleValidData(t *testing.T) {
	testCases := map[string]struct {
		userTitle domain.UserTitle
	}{
		"Can empty validations errors": {
			userTitle: domain.UserTitle{
				ID:        1,
				Enabled:   false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Title: domain.Title{
					Title:       "Title 1",
					Description: "Description 1",
					Purchasable: false,
					Status:      "available",
				},
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

			err := tc.userTitle.ValidateUserTitle()

			assert.NoError(t, err)
		})
	}
}

func TestCreateUserTitleWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		title   domain.UserTitle
		wantErr string
	}{
		"Missing required fields": {
			title: domain.UserTitle{},
			wantErr: regexp.QuoteMeta(`
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
			`),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.title.ValidateUserTitle()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
