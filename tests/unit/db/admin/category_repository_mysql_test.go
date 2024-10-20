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

func TestAdminCategoryRepositoryMySQL_GetAll(t *testing.T) {
	fixedTime := time.Now()
	gormDB, mock := testutils.Setup(t)

	repo := db_admin.NewAdminCategoryRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		mockBehavior func()
		expectedLen  int
		expectedErr  error
	}{
		"success case": {
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(1, "Category 1", fixedTime, fixedTime).
					AddRow(2, "Category 2", fixedTime, fixedTime)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`deleted_at` IS NULL")).
					WillReturnRows(rows)
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		"no records found": {
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`deleted_at` IS NULL")).
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

func TestAdminCategoryRepositoryMySQL_Create(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		category     *domain.Category
		mockBehavior func(mock sqlmock.Sqlmock, category *domain.Category)
		expectedErr  error
		expectedSlug string
	}{
		"success case": {
			category: &domain.Category{
				Name:      "Category 1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, category *domain.Category) {
				expectedSlug := utils.Slugify(category.Name)

				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `categories`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						category.Name,
						expectedSlug,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr:  nil,
			expectedSlug: utils.Slugify("Category 1"),
		},
		"Failure - Insert Error": {
			category: &domain.Category{
				Name:      "Category 1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, category *domain.Category) {
				expectedSlug := utils.Slugify(category.Name)
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `categories`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						category.Name,
						expectedSlug,
					).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectedErr:  fmt.Errorf("database error"),
			expectedSlug: utils.Slugify("Category 1"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gormDB, mock := testutils.Setup(t)

			repo := db_admin.NewAdminCategoryRepositoryMySQL(gormDB)

			tc.category.Slug = utils.Slugify(tc.category.Name)

			tc.mockBehavior(mock, tc.category)

			err := repo.Create(tc.category)

			assert.Equal(t, tc.expectedSlug, tc.category.Slug)
			assert.Equal(t, tc.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdatePicture(t *testing.T) {
	gormDB, mock := testutils.Setup(t)

	repo := db_admin.NewAdminCategoryRepositoryMySQL(gormDB)

	tests := map[string]struct {
		categoryID uint
		request    ports_admin.UpdateCategoryInterface
		mock       func(request ports_admin.UpdateCategoryInterface)
		expectErr  bool
	}{
		"successful picture update": {
			categoryID: 1,
			request: ports_admin.UpdateCategoryInterface{
				Name: "Category 1",
				Slug: "category-1",
			},
			mock: func(request ports_admin.UpdateCategoryInterface) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `categories` SET `name`=?,`slug`=?,`updated_at`=? WHERE id = ? AND `categories`.`deleted_at` IS NULL")).
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
			categoryID: 2,
			request: ports_admin.UpdateCategoryInterface{
				Name: "Category 1",
				Slug: "category-1",
			},
			mock: func(request ports_admin.UpdateCategoryInterface) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `categories` SET `name`=?,`slug`=?,`updated_at`=? WHERE id = ? AND `categories`.`deleted_at` IS NULL")).
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
			err := repo.Update(tt.categoryID, tt.request)

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAdminCategoryRepositoryMySQL_Delete(t *testing.T) {
	gormDB, mock := testutils.Setup(t)
	repo := db_admin.NewAdminCategoryRepositoryMySQL(gormDB)

	testCases := map[string]struct {
		categoryID   uint
		mockBehavior func()
		expectedErr  error
	}{
		"successful delete": {
			categoryID: 1,
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `categories` SET `deleted_at`=? WHERE `categories`.`id` = ? AND `categories`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		"category not found": {
			categoryID: 99,
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `categories` SET `deleted_at`=? WHERE `categories`.`id` = ? AND `categories`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), 99).
					WillReturnError(errors.NewHttpError(http.StatusNotFound, "category not found"))
				mock.ExpectRollback()
			},
			expectedErr: errors.NewHttpError(http.StatusNotFound, "category not found."),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			tc.mockBehavior()

			err := repo.Delete(tc.categoryID)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
