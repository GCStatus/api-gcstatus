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

func TestCreateCategoriable(t *testing.T) {
	testCases := map[string]struct {
		categoriable domain.Categoriable
		mockBehavior func(mock sqlmock.Sqlmock, categoriable domain.Categoriable)
		expectError  bool
	}{
		"Success": {
			categoriable: domain.Categoriable{
				CategoriableID:   1,
				CategoriableType: "games",
				CategoryID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, categoriable domain.Categoriable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `categoriables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						categoriable.CategoriableID,
						categoriable.CategoriableType,
						categoriable.CategoryID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			categoriable: domain.Categoriable{
				CategoriableID:   1,
				CategoriableType: "games",
				CategoryID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, categoriable domain.Categoriable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `categoriables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						categoriable.CategoriableID,
						categoriable.CategoriableType,
						categoriable.CategoryID,
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

			tc.mockBehavior(mock, tc.categoriable)

			err := db.Create(&tc.categoriable).Error

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

func TestUpdateCategoriable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		categoriable domain.Categoriable
		mockBehavior func(mock sqlmock.Sqlmock, categoriable domain.Categoriable)
		expectError  bool
	}{
		"Success": {
			categoriable: domain.Categoriable{
				ID:               1,
				CategoriableID:   1,
				CategoriableType: "games",
				CategoryID:       1,
				CreatedAt:        fixedTime,
				UpdatedAt:        fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, categoriable domain.Categoriable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `categoriables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						categoriable.CategoriableID,
						categoriable.CategoriableType,
						categoriable.CategoryID,
						categoriable.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			categoriable: domain.Categoriable{
				ID:               1,
				CategoriableID:   1,
				CategoriableType: "games",
				CategoryID:       1,
				CreatedAt:        fixedTime,
				UpdatedAt:        fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, categoriable domain.Categoriable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `categoriables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						categoriable.CategoriableID,
						categoriable.CategoriableType,
						categoriable.CategoryID,
						categoriable.ID,
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

			tc.mockBehavior(mock, tc.categoriable)

			err := db.Save(&tc.categoriable).Error

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

func TestSoftDeleteCategoriable(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		categoriableID uint
		mockBehavior   func(mock sqlmock.Sqlmock, categoriableID uint)
		wantErr        bool
	}{
		"Can soft delete a Categoriable": {
			categoriableID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, categoriableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `categoriables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), categoriableID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			categoriableID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, categoriableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `categoriables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete categoriable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.categoriableID)

			err := db.Delete(&domain.Categoriable{}, tc.categoriableID).Error

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
