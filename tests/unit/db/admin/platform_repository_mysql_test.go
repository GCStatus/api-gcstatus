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

func TestAdminPlatformRepositoryMySQL_GetAll(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := testutils.Setup(t)

	repo := db_admin.NewAdminPlatformRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		mockBehavior func()
		expectedLen  int
		expectedErr  error
	}{
		"success case": {
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(1, "Platform 1", fixedTime, fixedTime).
					AddRow(2, "Platform 2", fixedTime, fixedTime)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `platforms` WHERE `platforms`.`deleted_at` IS NULL")).
					WillReturnRows(rows)
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		"no records found": {
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `platforms` WHERE `platforms`.`deleted_at` IS NULL")).
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

func TestAdminPlatformRepositoryMySQL_Create(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		platform     *domain.Platform
		mockBehavior func(mock sqlmock.Sqlmock, platform *domain.Platform)
		expectedErr  error
		expectedSlug string
	}{
		"success case": {
			platform: &domain.Platform{
				Name:      "Platform 1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, platform *domain.Platform) {
				expectedSlug := utils.Slugify(platform.Name)

				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `platforms`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						platform.Name,
						expectedSlug,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr:  nil,
			expectedSlug: utils.Slugify("Platform 1"),
		},
		"Failure - Insert Error": {
			platform: &domain.Platform{
				Name:      "Platform 1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, platform *domain.Platform) {
				expectedSlug := utils.Slugify(platform.Name)
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `platforms`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						platform.Name,
						expectedSlug,
					).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectedErr:  fmt.Errorf("database error"),
			expectedSlug: utils.Slugify("Platform 1"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db_admin.NewAdminPlatformRepositoryMySQL(gormDB)

			tc.platform.Slug = utils.Slugify(tc.platform.Name)

			tc.mockBehavior(mock, tc.platform)

			err := repo.Create(tc.platform)

			assert.Equal(t, tc.expectedSlug, tc.platform.Slug)
			assert.Equal(t, tc.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAdminPlatformRepositoryMySQL_Update(t *testing.T) {
	gormDB, mock := testutils.Setup(t)

	repo := db_admin.NewAdminPlatformRepositoryMySQL(gormDB)

	tests := map[string]struct {
		platformID uint
		request    ports_admin.UpdatePlatformInterface
		mock       func(request ports_admin.UpdatePlatformInterface)
		expectErr  bool
	}{
		"successful picture update": {
			platformID: 1,
			request: ports_admin.UpdatePlatformInterface{
				Name: "Platform 1",
				Slug: "platform-1",
			},
			mock: func(request ports_admin.UpdatePlatformInterface) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `platforms` SET `name`=?,`slug`=?,`updated_at`=? WHERE id = ? AND `platforms`.`deleted_at` IS NULL")).
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
			platformID: 2,
			request: ports_admin.UpdatePlatformInterface{
				Name: "Platform 1",
				Slug: "platform-1",
			},
			mock: func(request ports_admin.UpdatePlatformInterface) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `platforms` SET `name`=?,`slug`=?,`updated_at`=? WHERE id = ? AND `platforms`.`deleted_at` IS NULL")).
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
			err := repo.Update(tt.platformID, tt.request)

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAdminPlatformRepositoryMySQL_Delete(t *testing.T) {
	gormDB, mock := testutils.Setup(t)
	repo := db_admin.NewAdminPlatformRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		platformID   uint
		mockBehavior func()
		expectedErr  error
	}{
		"successful delete": {
			platformID: 1,
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `platforms` SET `deleted_at`=? WHERE `platforms`.`id` = ? AND `platforms`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		"platform not found": {
			platformID: 99,
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `platforms` SET `deleted_at`=? WHERE `platforms`.`id` = ? AND `platforms`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), 99).
					WillReturnError(errors.NewHttpError(http.StatusNotFound, "platform not found"))
				mock.ExpectRollback()
			},
			expectedErr: errors.NewHttpError(http.StatusNotFound, "platform not found"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			tc.mockBehavior()

			err := repo.Delete(tc.platformID)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
