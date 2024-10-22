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

func TestAdminGenreRepositoryMySQL_GetAll(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := testutils.Setup(t)

	repo := db_admin.NewAdminGenreRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		mockBehavior func()
		expectedLen  int
		expectedErr  error
	}{
		"success case": {
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(1, "Genre 1", fixedTime, fixedTime).
					AddRow(2, "Genre 2", fixedTime, fixedTime)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `genres` WHERE `genres`.`deleted_at` IS NULL")).
					WillReturnRows(rows)
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		"no records found": {
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `genres` WHERE `genres`.`deleted_at` IS NULL")).
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

func TestAdminGenreRepositoryMySQL_Create(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		genre        *domain.Genre
		mockBehavior func(mock sqlmock.Sqlmock, genre *domain.Genre)
		expectedErr  error
		expectedSlug string
	}{
		"success case": {
			genre: &domain.Genre{
				Name:      "Genre 1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, genre *domain.Genre) {
				expectedSlug := utils.Slugify(genre.Name)

				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `genres`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						genre.Name,
						expectedSlug,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr:  nil,
			expectedSlug: utils.Slugify("Genre 1"),
		},
		"Failure - Insert Error": {
			genre: &domain.Genre{
				Name:      "Genre 1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, genre *domain.Genre) {
				expectedSlug := utils.Slugify(genre.Name)
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `genres`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						genre.Name,
						expectedSlug,
					).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectedErr:  fmt.Errorf("database error"),
			expectedSlug: utils.Slugify("Genre 1"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db_admin.NewAdminGenreRepositoryMySQL(gormDB)

			tc.genre.Slug = utils.Slugify(tc.genre.Name)

			tc.mockBehavior(mock, tc.genre)

			err := repo.Create(tc.genre)

			assert.Equal(t, tc.expectedSlug, tc.genre.Slug)
			assert.Equal(t, tc.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAdminGenreRepositoryMySQL_Update(t *testing.T) {
	gormDB, mock := testutils.Setup(t)

	repo := db_admin.NewAdminGenreRepositoryMySQL(gormDB)

	tests := map[string]struct {
		genreID   uint
		request   ports_admin.UpdateGenreInterface
		mock      func(request ports_admin.UpdateGenreInterface)
		expectErr bool
	}{
		"successful picture update": {
			genreID: 1,
			request: ports_admin.UpdateGenreInterface{
				Name: "Genre 1",
				Slug: "genre-1",
			},
			mock: func(request ports_admin.UpdateGenreInterface) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `genres` SET `name`=?,`slug`=?,`updated_at`=? WHERE id = ? AND `genres`.`deleted_at` IS NULL")).
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
			genreID: 2,
			request: ports_admin.UpdateGenreInterface{
				Name: "Genre 1",
				Slug: "genre-1",
			},
			mock: func(request ports_admin.UpdateGenreInterface) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `genres` SET `name`=?,`slug`=?,`updated_at`=? WHERE id = ? AND `genres`.`deleted_at` IS NULL")).
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
			err := repo.Update(tt.genreID, tt.request)

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAdminGenreRepositoryMySQL_Delete(t *testing.T) {
	gormDB, mock := testutils.Setup(t)
	repo := db_admin.NewAdminGenreRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		genreID      uint
		mockBehavior func()
		expectedErr  error
	}{
		"successful delete": {
			genreID: 1,
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `genres` SET `deleted_at`=? WHERE `genres`.`id` = ? AND `genres`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		"genre not found": {
			genreID: 99,
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `genres` SET `deleted_at`=? WHERE `genres`.`id` = ? AND `genres`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), 99).
					WillReturnError(errors.NewHttpError(http.StatusNotFound, "genre not found"))
				mock.ExpectRollback()
			},
			expectedErr: errors.NewHttpError(http.StatusNotFound, "genre not found"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			tc.mockBehavior()

			err := repo.Delete(tc.genreID)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
