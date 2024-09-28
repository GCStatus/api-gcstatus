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

func CreateTitleTest(t *testing.T) {
	testCases := map[string]struct {
		title        domain.Title
		mockBehavior func(mock sqlmock.Sqlmock, title domain.Title)
		expectErr    bool
	}{
		"Successfully created": {
			title: domain.Title{
				Title:       "Title 1",
				Description: "Title 1",
				Purchasable: false,
				Status:      "available",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, title domain.Title) {
				mock.ExpectBegin()

				mock.ExpectExec("^INSERT INTO `titles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						title.Title,
						title.Description,
						title.Purchasable,
						title.Status,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		"Failure - Insert Error": {
			title: domain.Title{
				Title:       "Title 1",
				Description: "Title 1",
				Purchasable: false,
				Status:      "available",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, title domain.Title) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `titles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						title.Title,
						title.Description,
						title.Purchasable,
						title.Status,
					).
					WillReturnError(fmt.Errorf("some error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := tests.Setup(t)

			tc.mockBehavior(mock, tc.title)

			err := db.Create(&tc.title).Error

			if tc.expectErr {
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

func TestSoftDeleteTitle(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		titleID      uint
		mockBehavior func(mock sqlmock.Sqlmock, titleID uint)
		wantErr      bool
	}{
		"Can soft delete a title": {
			titleID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, titleID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `titles` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), titleID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			titleID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, titleID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `titles` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete title"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.titleID)

			err := db.Delete(&domain.Title{}, tc.titleID).Error

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

func TestUpdateTitle(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		title        domain.Title
		mockBehavior func(mock sqlmock.Sqlmock, title domain.Title)
		expectError  bool
	}{
		"Success": {
			title: domain.Title{
				ID:          1,
				Title:       "Title 1",
				Description: "Title 1",
				Purchasable: false,
				Status:      "available",
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, title domain.Title) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `titles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						title.Title,
						title.Description,
						title.Cost,
						title.Purchasable,
						title.Status,
						title.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			title: domain.Title{
				ID:          1,
				Title:       "Title 1",
				Description: "Title 1",
				Purchasable: false,
				Status:      "available",
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, title domain.Title) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `titles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						title.Title,
						title.Description,
						title.Cost,
						title.Purchasable,
						title.Status,
						title.ID,
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

			tc.mockBehavior(mock, tc.title)

			err := db.Save(&tc.title).Error

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

func TestValidateTitleValidData(t *testing.T) {
	testCases := map[string]struct {
		title domain.Title
	}{
		"Can empty validations errors": {
			title: domain.Title{
				ID:          1,
				Title:       "Title 1",
				Description: "Title 1",
				Purchasable: false,
				Status:      "available",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.title.ValidateTitle()

			assert.NoError(t, err)
		})
	}
}

func TestCreateTitleWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		title   domain.Title
		wantErr string
	}{
		"Missing required fields": {
			title:   domain.Title{},
			wantErr: "Title is a required field, Description is a required field, Status is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.title.ValidateTitle()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
