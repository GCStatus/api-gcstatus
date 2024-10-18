package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTag(t *testing.T) {
	testCases := map[string]struct {
		tag          domain.Tag
		mockBehavior func(mock sqlmock.Sqlmock, tag domain.Tag)
		expectError  bool
	}{
		"Success": {
			tag: domain.Tag{
				Name: "Tag 1",
				Slug: "tag-1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, tag domain.Tag) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tags`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						tag.Name,
						tag.Slug,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			tag: domain.Tag{
				Name: "Failure",
				Slug: "failure",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, tag domain.Tag) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tags`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						tag.Name,
						tag.Slug,
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

			tc.mockBehavior(mock, tc.tag)

			err := db.Create(&tc.tag).Error

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

func TestUpdateTag(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		tag          domain.Tag
		mockBehavior func(mock sqlmock.Sqlmock, tag domain.Tag)
		expectError  bool
	}{
		"Success": {
			tag: domain.Tag{
				ID:        1,
				Name:      "Tag 1",
				Slug:      "tag-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, tag domain.Tag) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tags`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						tag.Name,
						tag.Slug,
						tag.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			tag: domain.Tag{
				ID:        1,
				Name:      "Tag 1",
				Slug:      "tag-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, tag domain.Tag) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tags`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						tag.Name,
						tag.Slug,
						tag.ID,
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

			tc.mockBehavior(mock, tc.tag)

			err := db.Save(&tc.tag).Error

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

func TestSoftDeleteTag(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		tagID        uint
		mockBehavior func(mock sqlmock.Sqlmock, tagID uint)
		wantErr      bool
	}{
		"Can soft delete a Tag": {
			tagID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, tagID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tags` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), tagID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			tagID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, tagID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `tags` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Tag"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.tagID)

			err := db.Delete(&domain.Tag{}, tc.tagID).Error

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

func TestGetTagByID(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		tagID     uint
		mockFunc  func()
		wantTag   domain.Tag
		wantError bool
	}{
		"Valid Tag fetch": {
			tagID: 1,
			wantTag: domain.Tag{
				ID:   1,
				Name: "Tag 1",
			},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Tag 1")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tags` WHERE `tags`.`id` = ? AND `tags`.`deleted_at` IS NULL ORDER BY `tags`.`id` LIMIT ?")).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			wantError: false,
		},
		"Tag not found": {
			tagID:     2,
			wantTag:   domain.Tag{},
			wantError: true,
			mockFunc: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tags` WHERE `tags`.`id` = ? AND `tags`.`deleted_at` IS NULL ORDER BY `tags`.`id` LIMIT ?")).
					WithArgs(2, 1).WillReturnError(fmt.Errorf("record not found"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			var tag domain.Tag
			err := db.First(&tag, tc.tagID).Error

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantTag, tag)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidateTagValidData(t *testing.T) {
	testCases := map[string]struct {
		tag domain.Tag
	}{
		"Can empty validations errors": {
			tag: domain.Tag{
				Name: "Tag 1",
				Slug: "tag-1",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.tag.ValidateTag()
			assert.NoError(t, err)
		})
	}
}

func TestCreateTagWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		tag     domain.Tag
		wantErr string
	}{
		"Missing required fields": {
			tag:     domain.Tag{},
			wantErr: "Name is a required field, Slug is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.tag.ValidateTag()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
