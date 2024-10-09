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

func TestCreateTaggable(t *testing.T) {
	testCases := map[string]struct {
		taggable     domain.Taggable
		mockBehavior func(mock sqlmock.Sqlmock, taggable domain.Taggable)
		expectError  bool
	}{
		"Success": {
			taggable: domain.Taggable{
				TaggableID:   1,
				TaggableType: "games",
				TagID:        1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, taggable domain.Taggable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `taggables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						taggable.TaggableID,
						taggable.TaggableType,
						taggable.TagID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			taggable: domain.Taggable{
				TaggableID:   1,
				TaggableType: "games",
				TagID:        1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, taggable domain.Taggable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `taggables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						taggable.TaggableID,
						taggable.TaggableType,
						taggable.TagID,
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

			tc.mockBehavior(mock, tc.taggable)

			err := db.Create(&tc.taggable).Error

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

func TestUpdateTaggable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		taggable     domain.Taggable
		mockBehavior func(mock sqlmock.Sqlmock, taggable domain.Taggable)
		expectError  bool
	}{
		"Success": {
			taggable: domain.Taggable{
				ID:           1,
				TaggableID:   1,
				TaggableType: "games",
				TagID:        1,
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, taggable domain.Taggable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `taggables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						taggable.TaggableID,
						taggable.TaggableType,
						taggable.TagID,
						taggable.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			taggable: domain.Taggable{
				ID:           1,
				TaggableID:   1,
				TaggableType: "games",
				TagID:        1,
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, taggable domain.Taggable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `taggables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						taggable.TaggableID,
						taggable.TaggableType,
						taggable.TagID,
						taggable.ID,
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

			tc.mockBehavior(mock, tc.taggable)

			err := db.Save(&tc.taggable).Error

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

func TestSoftDeleteTaggable(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		taggableID   uint
		mockBehavior func(mock sqlmock.Sqlmock, taggableID uint)
		wantErr      bool
	}{
		"Can soft delete a Taggable": {
			taggableID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, taggableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `taggables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), taggableID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			taggableID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, taggableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `taggables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Taggable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.taggableID)

			err := db.Delete(&domain.Taggable{}, tc.taggableID).Error

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
