package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateViewable(t *testing.T) {
	testCases := map[string]struct {
		viewable     domain.Viewable
		mockBehavior func(mock sqlmock.Sqlmock, viewable domain.Viewable)
		expectError  bool
	}{
		"Success": {
			viewable: domain.Viewable{
				ViewableID:   1,
				ViewableType: "games",
				UserID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, viewable domain.Viewable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `viewables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						viewable.ViewableID,
						viewable.ViewableType,
						viewable.UserID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			viewable: domain.Viewable{
				ViewableID:   1,
				ViewableType: "games",
				UserID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, viewable domain.Viewable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `viewables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						viewable.ViewableID,
						viewable.ViewableType,
						viewable.UserID,
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

			tc.mockBehavior(mock, tc.viewable)

			err := db.Create(&tc.viewable).Error

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

func TestUpdateViewable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		viewable     domain.Viewable
		mockBehavior func(mock sqlmock.Sqlmock, viewable domain.Viewable)
		expectError  bool
	}{
		"Success": {
			viewable: domain.Viewable{
				ID:           1,
				ViewableID:   1,
				ViewableType: "games",
				UserID:       1,
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, viewable domain.Viewable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `viewables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						viewable.ViewableID,
						viewable.ViewableType,
						viewable.UserID,
						viewable.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			viewable: domain.Viewable{
				ID:           1,
				ViewableID:   1,
				ViewableType: "games",
				UserID:       1,
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, viewable domain.Viewable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `viewables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						viewable.ViewableID,
						viewable.ViewableType,
						viewable.UserID,
						viewable.ID,
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

			tc.mockBehavior(mock, tc.viewable)

			err := db.Save(&tc.viewable).Error

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

func TestSoftDeleteViewable(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		viewableID   uint
		mockBehavior func(mock sqlmock.Sqlmock, viewableID uint)
		wantErr      bool
	}{
		"Can soft delete a Viewable": {
			viewableID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, viewableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `viewables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), viewableID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			viewableID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, viewableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `viewables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Viewable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.viewableID)

			err := db.Delete(&domain.Viewable{}, tc.viewableID).Error

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

func TestValidateViewableValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		viewable domain.Viewable
	}{
		"Valid Viewable with zero amount": {
			viewable: domain.Viewable{
				ID:           1,
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
				ViewableID:   1,
				ViewableType: "games",
				User: domain.User{
					Name:       "John Doe",
					Email:      "johndoe@example.com",
					Nickname:   "johnny",
					Blocked:    false,
					Experience: 500,
					Birthdate:  fixedTime,
					Password:   "supersecretpassword",
					Profile: domain.Profile{
						Share: true,
					},
					Wallet: domain.Wallet{
						Amount: 10,
					},
					Level: domain.Level{
						Level:      1,
						Experience: 500,
						Coins:      10,
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.viewable.ValidateViewable()
			assert.NoError(t, err)
		})
	}
}
