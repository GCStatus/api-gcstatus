package tests

import (
	"errors"
	"fmt"
	"gcstatus/internal/adapters/db"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCommentRepositoryMySQL_FindById(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := testutils.Setup(t)
	repo := db.NewCommentRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		commentID       uint
		mockSetup       func()
		expectedError   error
		expectedComment *domain.Commentable
	}{
		"success - comment found": {
			commentID: 1,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `commentables` WHERE `commentables`.`id` = ? AND `commentables`.`deleted_at` IS NULL ORDER BY `commentables`.`id` LIMIT ?")).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "comment", "commentable_id", "commentable_type", "user_id", "created_at", "updated_at"}).
						AddRow(1, "Base comment.", 1, "games", 1, fixedTime, fixedTime))
			},
			expectedComment: &domain.Commentable{
				ID:              1,
				Comment:         "Base comment.",
				CommentableID:   1,
				CommentableType: "games",
				UserID:          1,
				CreatedAt:       fixedTime,
				UpdatedAt:       fixedTime,
			},
			expectedError: nil,
		},
		"error - comment not found": {
			commentID: 2,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `commentables` WHERE `commentables`.`id` = ? AND `commentables`.`deleted_at` IS NULL ORDER BY `commentables`.`id` LIMIT ?")).
					WithArgs(2, 1).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectedComment: nil,
			expectedError:   gorm.ErrRecordNotFound,
		},
		"error - db failure": {
			commentID: 3,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `commentables` WHERE `commentables`.`id` = ? AND `commentables`.`deleted_at` IS NULL ORDER BY `commentables`.`id` LIMIT ?")).
					WithArgs(3, 1).
					WillReturnError(errors.New("db error"))
			},
			expectedComment: nil,
			expectedError:   errors.New("db error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()

			comment, err := repo.FindByID(tc.commentID)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedComment, comment)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCommentRepositoryMySQL_Create(t *testing.T) {
	testCases := map[string]struct {
		commentable  *domain.Commentable
		mockBehavior func(mock sqlmock.Sqlmock, commentable *domain.Commentable)
		expectedErr  error
	}{
		"success case": {
			commentable: &domain.Commentable{
				Comment:         "Base comment.",
				CommentableID:   1,
				CommentableType: "games",
				UserID:          1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, commentable *domain.Commentable) {
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

				commentRows := sqlmock.NewRows([]string{"id", "user_id", "commentable_id", "commentable_type", "comment", "parent_id"}).
					AddRow(1, 1, 1, "games", "Base comment.", nil)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `commentables` WHERE `commentables`.`id` = ? AND `commentables`.`deleted_at` IS NULL AND `commentables`.`id` = ? ORDER BY `commentables`.`id` LIMIT ?")).
					WithArgs(1, 1, 1).WillReturnRows(commentRows)

				heartsRows := mock.NewRows([]string{"id", "heartable_id", "heartable_type", "user_id"}).
					AddRow(1, 1, "commentables", 1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `heartables` WHERE `heartable_type` = ? AND `heartables`.`heartable_id` = ? AND `heartables`.`deleted_at` IS NULL")).
					WithArgs("commentables", 1).
					WillReturnRows(heartsRows)

				repliesRows := sqlmock.NewRows([]string{"id", "user_id", "commentable_id", "commentable_type", "comment", "parent_id"})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `commentables` WHERE `commentables`.`parent_id` = ? AND `commentables`.`deleted_at` IS NULL")).
					WithArgs(1).WillReturnRows(repliesRows)

				userRows := sqlmock.NewRows([]string{"id", "name", "email", "nickname", "created_at", "updated_at"}).
					AddRow(1, "Fake", "fake@gmail.com", "fake", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")).
					WithArgs(1).WillReturnRows(userRows)
			},
			expectedErr: nil,
		},
		"Failure - Insert Error": {
			commentable: &domain.Commentable{
				Comment:         "Base comment.",
				CommentableID:   1,
				CommentableType: "games",
				UserID:          1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, commentable *domain.Commentable) {
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
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectedErr: fmt.Errorf("database error"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewCommentRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.commentable)

			_, err := repo.Create(*tc.commentable)

			assert.Equal(t, tc.expectedErr, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCommentRepositoryMySQL_Delete(t *testing.T) {
	testCases := map[string]struct {
		commentableID uint
		mockBehavior  func(mock sqlmock.Sqlmock, commentableID uint)
		wantErr       bool
	}{
		"can delete a commentable": {
			commentableID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, commentableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `commentables` SET `deleted_at`=? WHERE `commentables`.`id` = ? AND `commentables`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), commentableID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"delete fails": {
			commentableID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, commentableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `commentables` SET `deleted_at`=? WHERE `commentables`.`id` = ? AND `commentables`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete commentable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db.NewCommentRepositoryMySQL(gormDB)

			tc.mockBehavior(mock, tc.commentableID)

			err := repo.Delete(tc.commentableID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete commentable")
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
