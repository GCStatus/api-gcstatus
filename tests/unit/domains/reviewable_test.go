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

func TestCreateReviewable(t *testing.T) {
	testCases := map[string]struct {
		reviewable   domain.Reviewable
		mockBehavior func(mock sqlmock.Sqlmock, reviewable domain.Reviewable)
		expectError  bool
	}{
		"Success": {
			reviewable: domain.Reviewable{
				Rate:           5,
				Review:         "Good game!",
				Played:         true,
				ReviewableID:   1,
				ReviewableType: "games",
				UserID:         1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, reviewable domain.Reviewable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `reviewables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						reviewable.Rate,
						reviewable.Review,
						reviewable.Played,
						reviewable.ReviewableID,
						reviewable.ReviewableType,
						reviewable.UserID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			reviewable: domain.Reviewable{
				Rate:           5,
				Review:         "Good game!",
				Played:         true,
				ReviewableID:   1,
				ReviewableType: "games",
				UserID:         1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, reviewable domain.Reviewable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `reviewables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						reviewable.Rate,
						reviewable.Review,
						reviewable.Played,
						reviewable.ReviewableID,
						reviewable.ReviewableType,
						reviewable.UserID,
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

			tc.mockBehavior(mock, tc.reviewable)

			err := db.Create(&tc.reviewable).Error

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

func TestUpdateReviewable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		reviewable   domain.Reviewable
		mockBehavior func(mock sqlmock.Sqlmock, reviewable domain.Reviewable)
		expectError  bool
	}{
		"Success": {
			reviewable: domain.Reviewable{
				ID:             1,
				Rate:           5,
				Review:         "Good game!",
				Played:         true,
				ReviewableID:   1,
				ReviewableType: "games",
				UserID:         1,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, reviewable domain.Reviewable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `reviewables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						reviewable.Rate,
						reviewable.Review,
						reviewable.Played,
						reviewable.ReviewableID,
						reviewable.ReviewableType,
						reviewable.UserID,
						reviewable.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			reviewable: domain.Reviewable{
				ID:             1,
				Rate:           5,
				Review:         "Good game!",
				Played:         true,
				ReviewableID:   1,
				ReviewableType: "games",
				UserID:         1,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, reviewable domain.Reviewable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `reviewables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						reviewable.Rate,
						reviewable.Review,
						reviewable.Played,
						reviewable.ReviewableID,
						reviewable.ReviewableType,
						reviewable.UserID,
						reviewable.ID,
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

			tc.mockBehavior(mock, tc.reviewable)

			err := db.Save(&tc.reviewable).Error

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

func TestSoftDeleteReviewable(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		reviewableID uint
		mockBehavior func(mock sqlmock.Sqlmock, reviewableID uint)
		wantErr      bool
	}{
		"Can soft delete a Reviewable": {
			reviewableID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, reviewableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `reviewables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), reviewableID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			reviewableID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, reviewableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `reviewables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Reviewable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.reviewableID)

			err := db.Delete(&domain.Reviewable{}, tc.reviewableID).Error

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

func TestValidateReviewableValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		reviewable domain.Reviewable
	}{
		"Valid Reviewable with zero amount": {
			reviewable: domain.Reviewable{
				ID:             1,
				Rate:           5,
				Review:         "Good game!",
				Played:         true,
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
				ReviewableID:   1,
				ReviewableType: "games",
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
			err := tc.reviewable.ValidateReviewable()
			assert.NoError(t, err)
		})
	}
}

func TestCreateReviewableWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		reviewable domain.Reviewable
		wantErr    string
	}{
		"Missing required fields": {
			reviewable: domain.Reviewable{},
			wantErr: `
				Rate is a required field,
				Review is a required field,
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
			err := tc.reviewable.ValidateReviewable()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
