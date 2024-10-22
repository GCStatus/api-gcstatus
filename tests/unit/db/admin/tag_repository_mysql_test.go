package tests

import (
	"fmt"
	db_admin "gcstatus/internal/adapters/db/admin"
	"gcstatus/internal/domain"
	"gcstatus/internal/errors"
	ports_admin "gcstatus/internal/ports/admin"
	"gcstatus/internal/utils"
	testutils "gcstatus/tests/utils"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAdminTagRepositoryMySQL_GetAll(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := testutils.Setup(t)

	repo := db_admin.NewAdminTagRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		mockBehavior func()
		expectedLen  int
		expectedErr  error
	}{
		"success case": {
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(1, "Tag 1", fixedTime, fixedTime).
					AddRow(2, "Tag 2", fixedTime, fixedTime)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tags` WHERE `tags`.`deleted_at` IS NULL")).
					WillReturnRows(rows)
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		"no records found": {
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tags` WHERE `tags`.`deleted_at` IS NULL")).
					WillReturnRows(rows)
			},
			expectedLen: 0,
			expectedErr: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior()

			levels, err := repo.GetAll()

			assert.Equal(t, tc.expectedErr, err)
			assert.Len(t, levels, tc.expectedLen)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAdminTagRepositoryMySQL_Create(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		tag          *domain.Tag
		mockBehavior func(mock sqlmock.Sqlmock, tag *domain.Tag)
		expectedErr  error
		expectedSlug string
	}{
		"success case": {
			tag: &domain.Tag{
				Name:      "Tag 1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, tag *domain.Tag) {
				expectedSlug := utils.Slugify(tag.Name)

				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tags`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						tag.Name,
						expectedSlug,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr:  nil,
			expectedSlug: utils.Slugify("Tag 1"),
		},
		"Failure - Insert Error": {
			tag: &domain.Tag{
				Name:      "Tag 1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, tag *domain.Tag) {
				expectedSlug := utils.Slugify(tag.Name)
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `tags`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						tag.Name,
						expectedSlug,
					).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectedErr:  fmt.Errorf("database error"),
			expectedSlug: utils.Slugify("Tag 1"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db_admin.NewAdminTagRepositoryMySQL(gormDB)

			tc.tag.Slug = utils.Slugify(tc.tag.Name)

			tc.mockBehavior(mock, tc.tag)

			err := repo.Create(tc.tag)

			assert.Equal(t, tc.expectedSlug, tc.tag.Slug)
			assert.Equal(t, tc.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAdminTagRepositoryMySQL_Update(t *testing.T) {
	gormDB, mock := testutils.Setup(t)

	repo := db_admin.NewAdminTagRepositoryMySQL(gormDB)

	tests := map[string]struct {
		tagID     uint
		request   ports_admin.UpdateTagInterface
		mock      func(request ports_admin.UpdateTagInterface)
		expectErr bool
	}{
		"successful picture update": {
			tagID: 1,
			request: ports_admin.UpdateTagInterface{
				Name: "Tag 1",
				Slug: "tag-1",
			},
			mock: func(request ports_admin.UpdateTagInterface) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `tags` SET `name`=?,`slug`=?,`updated_at`=? WHERE id = ? AND `tags`.`deleted_at` IS NULL")).
					WithArgs(
						request.Name,
						request.Slug,
						sqlmock.AnyArg(),
						1,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		"failed picture update due to database error": {
			tagID: 2,
			request: ports_admin.UpdateTagInterface{
				Name: "Tag 1",
				Slug: "tag-1",
			},
			mock: func(request ports_admin.UpdateTagInterface) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `tags` SET `name`=?,`slug`=?,`updated_at`=? WHERE id = ? AND `tags`.`deleted_at` IS NULL")).
					WithArgs(
						request.Name,
						request.Slug,
						sqlmock.AnyArg(),
						2,
					).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.mock(tt.request)
			err := repo.Update(tt.tagID, tt.request)

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAdminTagRepositoryMySQL_Delete(t *testing.T) {
	gormDB, mock := testutils.Setup(t)
	repo := db_admin.NewAdminTagRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		tagID        uint
		mockBehavior func()
		expectedErr  error
	}{
		"successful delete": {
			tagID: 1,
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `tags` SET `deleted_at`=? WHERE `tags`.`id` = ? AND `tags`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		"tag not found": {
			tagID: 99,
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `tags` SET `deleted_at`=? WHERE `tags`.`id` = ? AND `tags`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), 99).
					WillReturnError(errors.NewHttpError(http.StatusNotFound, "tag not found"))
				mock.ExpectRollback()
			},
			expectedErr: errors.NewHttpError(http.StatusNotFound, "tag not found"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			tc.mockBehavior()

			err := repo.Delete(tc.tagID)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
