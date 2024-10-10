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

func TestCreateGalleriable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		galleriable  domain.Galleriable
		mockBehavior func(mock sqlmock.Sqlmock, galleriable domain.Galleriable)
		expectError  bool
	}{
		"Success": {
			galleriable: domain.Galleriable{
				S3:              false,
				Path:            "https://google.com",
				GalleriableID:   1,
				GalleriableType: "games",
				MediaTypeID:     1,
				CreatedAt:       fixedTime,
				UpdatedAt:       fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, galleriable domain.Galleriable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `galleriables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						galleriable.S3,
						galleriable.Path,
						galleriable.GalleriableID,
						galleriable.GalleriableType,
						galleriable.MediaTypeID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			galleriable: domain.Galleriable{
				S3:              false,
				Path:            "https://google.com",
				GalleriableID:   1,
				GalleriableType: "games",
				MediaTypeID:     1,
				CreatedAt:       fixedTime,
				UpdatedAt:       fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, galleriable domain.Galleriable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `galleriables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						galleriable.S3,
						galleriable.Path,
						galleriable.GalleriableID,
						galleriable.GalleriableType,
						galleriable.MediaTypeID,
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

			tc.mockBehavior(mock, tc.galleriable)

			err := db.Create(&tc.galleriable).Error

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

func TestUpdateGalleriable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		galleriable  domain.Galleriable
		mockBehavior func(mock sqlmock.Sqlmock, galleriable domain.Galleriable)
		expectError  bool
	}{
		"Success": {
			galleriable: domain.Galleriable{
				ID:              1,
				S3:              false,
				Path:            "https://google.com",
				GalleriableID:   1,
				GalleriableType: "games",
				MediaTypeID:     1,
				CreatedAt:       fixedTime,
				UpdatedAt:       fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, galleriable domain.Galleriable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `galleriables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						galleriable.S3,
						galleriable.Path,
						galleriable.GalleriableID,
						galleriable.GalleriableType,
						galleriable.MediaTypeID,
						galleriable.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			galleriable: domain.Galleriable{
				ID:              1,
				S3:              false,
				Path:            "https://google.com",
				GalleriableID:   1,
				GalleriableType: "games",
				MediaTypeID:     1,
				CreatedAt:       fixedTime,
				UpdatedAt:       fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, galleriable domain.Galleriable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `galleriables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						galleriable.S3,
						galleriable.Path,
						galleriable.GalleriableID,
						galleriable.GalleriableType,
						galleriable.MediaTypeID,
						galleriable.ID,
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

			tc.mockBehavior(mock, tc.galleriable)

			err := db.Save(&tc.galleriable).Error

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

func TestSoftDeleteGalleriable(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		galleriableID uint
		mockBehavior  func(mock sqlmock.Sqlmock, galleriableID uint)
		wantErr       bool
	}{
		"Can soft delete a Galleriable": {
			galleriableID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, galleriableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `galleriables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), galleriableID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			galleriableID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, galleriableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `galleriables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Galleriable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.galleriableID)

			err := db.Delete(&domain.Galleriable{}, tc.galleriableID).Error

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

func TestValidateGalleriableValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		galleriable domain.Galleriable
	}{
		"Can empty validations errors": {
			galleriable: domain.Galleriable{
				S3:              false,
				Path:            "https://google.com",
				GalleriableID:   1,
				GalleriableType: "games",
				CreatedAt:       fixedTime,
				UpdatedAt:       fixedTime,
				MediaType: domain.MediaType{
					ID:   1,
					Name: "photo",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.galleriable.ValidateGalleriable()
			assert.NoError(t, err)
		})
	}
}

func TestCreateGalleriableWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		galleriable domain.Galleriable
		wantErr     string
	}{
		"Missing required fields": {
			galleriable: domain.Galleriable{},
			wantErr:     "Path is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.galleriable.ValidateGalleriable()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
