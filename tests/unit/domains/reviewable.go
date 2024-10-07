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

func TestCreateReviewable(t *testing.T) {
	testCases := map[string]struct {
		reviewable   domain.Reviewable
		mockBehavior func(mock sqlmock.Sqlmock, reviewable domain.Reviewable)
		expectError  bool
	}{
		"Success": {
			reviewable: domain.Reviewable{
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
			db, mock := tests.Setup(t)

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
			db, mock := tests.Setup(t)

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
	db, mock := tests.Setup(t)

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
