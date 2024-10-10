package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/tests"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategory(t *testing.T) {
	testCases := map[string]struct {
		category     domain.Category
		mockBehavior func(mock sqlmock.Sqlmock, category domain.Category)
		expectError  bool
	}{
		"Success": {
			category: domain.Category{
				Name: "Category 1",
				Slug: "category-1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, category domain.Category) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `categories`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						category.Name,
						category.Slug,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			category: domain.Category{
				Name: "Failure",
				Slug: "failure",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, category domain.Category) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `categories`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						category.Name,
						category.Slug,
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

			tc.mockBehavior(mock, tc.category)

			err := db.Create(&tc.category).Error

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

func TestUpdateCategory(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		category     domain.Category
		mockBehavior func(mock sqlmock.Sqlmock, category domain.Category)
		expectError  bool
	}{
		"Success": {
			category: domain.Category{
				ID:        1,
				Name:      "Category 1",
				Slug:      "category-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, category domain.Category) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `categories`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						category.Name,
						category.Slug,
						category.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			category: domain.Category{
				ID:        1,
				Name:      "Category 1",
				Slug:      "category-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, category domain.Category) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `categories`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						category.Name,
						category.Slug,
						category.ID,
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

			tc.mockBehavior(mock, tc.category)

			err := db.Save(&tc.category).Error

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

func TestSoftDeleteCategory(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		categoryID   uint
		mockBehavior func(mock sqlmock.Sqlmock, categoryID uint)
		wantErr      bool
	}{
		"Can soft delete a category": {
			categoryID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, categoryID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `categories` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), categoryID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			categoryID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, categoryID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `categories` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete category"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.categoryID)

			err := db.Delete(&domain.Category{}, tc.categoryID).Error

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

func TestGetCategoryByID(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		categoryID   uint
		mockFunc     func()
		wantCategory domain.Category
		wantError    bool
	}{
		"Valid category fetch": {
			categoryID: 1,
			wantCategory: domain.Category{
				ID:   1,
				Name: "Category 1",
			},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Category 1")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` = ? AND `categories`.`deleted_at` IS NULL ORDER BY `categories`.`id` LIMIT ?")).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			wantError: false,
		},
		"Category not found": {
			categoryID:   2,
			wantCategory: domain.Category{},
			wantError:    true,
			mockFunc: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `categories` WHERE `categories`.`id` = ? AND `categories`.`deleted_at` IS NULL ORDER BY `categories`.`id` LIMIT ?")).
					WithArgs(2, 1).WillReturnError(fmt.Errorf("record not found"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			var category domain.Category
			err := db.First(&category, tc.categoryID).Error

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantCategory, category)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidateCategoryValidData(t *testing.T) {
	testCases := map[string]struct {
		category domain.Category
	}{
		"Can empty validations errors": {
			category: domain.Category{
				Name: "Category 1",
				Slug: "category-1",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.category.ValidateCategory()
			assert.NoError(t, err)
		})
	}
}

func TestCreateCategoryWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		category domain.Category
		wantErr  string
	}{
		"Missing required fields": {
			category: domain.Category{},
			wantErr:  "Name is a required field, Slug is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.category.ValidateCategory()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
