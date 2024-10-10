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

func CreateMediaTypeTest(t *testing.T) {
	testCases := map[string]struct {
		mediaType    domain.MediaType
		mockBehavior func(mock sqlmock.Sqlmock, mediaType domain.MediaType)
		expectErr    bool
	}{
		"Successfully created": {
			mediaType: domain.MediaType{
				Name: "photo",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, mediaType domain.MediaType) {
				mock.ExpectBegin()

				mock.ExpectExec("^INSERT INTO `media_types`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						mediaType.Name,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		"Failure - Insert Error": {
			mediaType: domain.MediaType{
				Name: "photo",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, mediaType domain.MediaType) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `media_types`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						mediaType.Name,
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

			tc.mockBehavior(mock, tc.mediaType)

			err := db.Create(&tc.mediaType).Error

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

func TestSoftDeleteMediaType(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		mediaTypeID  uint
		mockBehavior func(mock sqlmock.Sqlmock, mediaTypeID uint)
		wantErr      bool
	}{
		"Can soft delete a MediaType": {
			mediaTypeID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, mediaTypeID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `media_types` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), mediaTypeID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			mediaTypeID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, mediaTypeID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `media_types` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete mediaType"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.mediaTypeID)

			err := db.Delete(&domain.MediaType{}, tc.mediaTypeID).Error

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

func TestUpdateMediaType(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		MediaType    domain.MediaType
		mockBehavior func(mock sqlmock.Sqlmock, MediaType domain.MediaType)
		expectError  bool
	}{
		"Success": {
			MediaType: domain.MediaType{
				ID:        1,
				Name:      "photo",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, mediaType domain.MediaType) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `media_types`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						mediaType.Name,
						mediaType.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			MediaType: domain.MediaType{
				ID:        1,
				Name:      "video",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, mediaType domain.MediaType) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `media_types`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						mediaType.Name,
						mediaType.ID,
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

			tc.mockBehavior(mock, tc.MediaType)

			err := db.Save(&tc.MediaType).Error

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

func TestValidateMediaTypeValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		mediaType domain.MediaType
	}{
		"Can empty validations errors": {
			mediaType: domain.MediaType{
				Name:      "photo",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.mediaType.ValidateMediaType()

			assert.NoError(t, err)
		})
	}
}

func TestCreateMediaTypeWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		mediaType domain.MediaType
		wantErr   string
	}{
		"Missing required fields": {
			mediaType: domain.MediaType{},
			wantErr:   "Name is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.mediaType.ValidateMediaType()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
