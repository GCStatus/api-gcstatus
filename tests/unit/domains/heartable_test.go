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

func TestCreateHeartable(t *testing.T) {
	testCases := map[string]struct {
		heartable    domain.Heartable
		mockBehavior func(mock sqlmock.Sqlmock, heartable domain.Heartable)
		expectError  bool
	}{
		"Success": {
			heartable: domain.Heartable{
				HeartableID:   1,
				HeartableType: "games",
				UserID:        1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, heartable domain.Heartable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `heartables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						heartable.HeartableID,
						heartable.HeartableType,
						heartable.UserID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			heartable: domain.Heartable{
				HeartableID:   1,
				HeartableType: "games",
				UserID:        1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, heartable domain.Heartable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `heartables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						heartable.HeartableID,
						heartable.HeartableType,
						heartable.UserID,
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

			tc.mockBehavior(mock, tc.heartable)

			err := db.Create(&tc.heartable).Error

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

func TestUpdateHeartable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		heartable    domain.Heartable
		mockBehavior func(mock sqlmock.Sqlmock, heartable domain.Heartable)
		expectError  bool
	}{
		"Success": {
			heartable: domain.Heartable{
				ID:            1,
				HeartableID:   1,
				HeartableType: "games",
				UserID:        1,
				CreatedAt:     fixedTime,
				UpdatedAt:     fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, heartable domain.Heartable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `heartables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						heartable.HeartableID,
						heartable.HeartableType,
						heartable.UserID,
						heartable.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			heartable: domain.Heartable{
				ID:            1,
				HeartableID:   1,
				HeartableType: "games",
				UserID:        1,
				CreatedAt:     fixedTime,
				UpdatedAt:     fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, heartable domain.Heartable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `heartables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						heartable.HeartableID,
						heartable.HeartableType,
						heartable.UserID,
						heartable.ID,
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

			tc.mockBehavior(mock, tc.heartable)

			err := db.Save(&tc.heartable).Error

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

func TestSoftDeleteHeartable(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		HeartableID  uint
		mockBehavior func(mock sqlmock.Sqlmock, HeartableID uint)
		wantErr      bool
	}{
		"Can soft delete a Heartable": {
			HeartableID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, HeartableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `heartables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), HeartableID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			HeartableID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, HeartableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `heartables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Heartable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.HeartableID)

			err := db.Delete(&domain.Heartable{}, tc.HeartableID).Error

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

func TestValidateHeartableValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		heartable domain.Heartable
	}{
		"Valid Heartable with zero amount": {
			heartable: domain.Heartable{
				ID:            1,
				CreatedAt:     fixedTime,
				UpdatedAt:     fixedTime,
				HeartableID:   1,
				HeartableType: "games",
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
			err := tc.heartable.ValidateHeartable()
			assert.NoError(t, err)
		})
	}
}

func TestCreateHeartableWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		heartable domain.Heartable
		wantErr   string
	}{
		"Missing required fields": {
			heartable: domain.Heartable{},
			wantErr:   "Share is a required field, Level is a required field, Experience is a required field, Coins is a required field, Amount is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.heartable.ValidateHeartable()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
