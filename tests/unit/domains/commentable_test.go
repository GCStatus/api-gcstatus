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

func TestCreateCommentable(t *testing.T) {
	testCases := map[string]struct {
		commentable  domain.Commentable
		mockBehavior func(mock sqlmock.Sqlmock, commentable domain.Commentable)
		expectError  bool
	}{
		"Success": {
			commentable: domain.Commentable{
				Comment:         "Comment",
				CommentableID:   1,
				CommentableType: "games",
				UserID:          1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, commentable domain.Commentable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `commentables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						commentable.Comment,
						commentable.UserID,
						commentable.CommentableID,
						commentable.CommentableType,
						commentable.ParentID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			commentable: domain.Commentable{
				Comment:         "Comment",
				CommentableID:   1,
				CommentableType: "games",
				UserID:          1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, commentable domain.Commentable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `commentables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						commentable.Comment,
						commentable.UserID,
						commentable.CommentableID,
						commentable.CommentableType,
						commentable.ParentID,
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

			tc.mockBehavior(mock, tc.commentable)

			err := db.Create(&tc.commentable).Error

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

func TestUpdateCommentable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		commentable  domain.Commentable
		mockBehavior func(mock sqlmock.Sqlmock, commentable domain.Commentable)
		expectError  bool
	}{
		"Success": {
			commentable: domain.Commentable{
				ID:              1,
				Comment:         "Comment",
				CommentableID:   1,
				CommentableType: "games",
				UserID:          1,
				CreatedAt:       fixedTime,
				UpdatedAt:       fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, commentable domain.Commentable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `commentables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						commentable.Comment,
						commentable.UserID,
						commentable.CommentableID,
						commentable.CommentableType,
						commentable.ParentID,
						commentable.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			commentable: domain.Commentable{
				ID:              1,
				Comment:         "Comment",
				CommentableID:   1,
				CommentableType: "games",
				UserID:          1,
				CreatedAt:       fixedTime,
				UpdatedAt:       fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, commentable domain.Commentable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `commentables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						commentable.Comment,
						commentable.UserID,
						commentable.CommentableID,
						commentable.CommentableType,
						commentable.ParentID,
						commentable.ID,
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

			tc.mockBehavior(mock, tc.commentable)

			err := db.Save(&tc.commentable).Error

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

func TestSoftDeleteCommentable(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		commentableID uint
		mockBehavior  func(mock sqlmock.Sqlmock, commentableID uint)
		wantErr       bool
	}{
		"Can soft delete a Commentable": {
			commentableID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, commentableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `commentables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), commentableID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			commentableID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, commentableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `commentables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete commentable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.commentableID)

			err := db.Delete(&domain.Commentable{}, tc.commentableID).Error

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

func TestValidateCommentableValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		commentable domain.Commentable
	}{
		"Can empty validations errors": {
			commentable: domain.Commentable{
				Comment:         "Comment",
				CommentableID:   1,
				CommentableType: "games",
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
			err := tc.commentable.ValidateCommentable()
			assert.NoError(t, err)
		})
	}
}

func TestCreateCommentableWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		commentable domain.Commentable
		wantErr     string
	}{
		"Missing required fields": {
			commentable: domain.Commentable{},
			wantErr:     "Name is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.commentable.ValidateCommentable()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
