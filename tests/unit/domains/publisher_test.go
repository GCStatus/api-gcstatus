package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	testutils "gcstatus/tests/utils"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePublisher(t *testing.T) {
	testCases := map[string]struct {
		publisher    domain.Publisher
		mockBehavior func(mock sqlmock.Sqlmock, publisher domain.Publisher)
		expectError  bool
	}{
		"Success": {
			publisher: domain.Publisher{
				Name:   "Game Science",
				Acting: false,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, publisher domain.Publisher) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `publishers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						publisher.Name,
						publisher.Acting,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			publisher: domain.Publisher{
				Name:   "Game Science",
				Acting: false,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, publisher domain.Publisher) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `publishers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						publisher.Name,
						publisher.Acting,
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

			tc.mockBehavior(mock, tc.publisher)

			err := db.Create(&tc.publisher).Error

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

func TestUpdatePublisher(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		publisher    domain.Publisher
		mockBehavior func(mock sqlmock.Sqlmock, publisher domain.Publisher)
		expectError  bool
	}{
		"Success": {
			publisher: domain.Publisher{
				ID:        1,
				Name:      "Game Science",
				Acting:    true,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, publisher domain.Publisher) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `publishers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						publisher.Name,
						publisher.Acting,
						publisher.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			publisher: domain.Publisher{
				ID:        1,
				Name:      "Game Science",
				Acting:    false,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, publisher domain.Publisher) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `publishers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						publisher.Name,
						publisher.Acting,
						publisher.ID,
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

			tc.mockBehavior(mock, tc.publisher)

			err := db.Save(&tc.publisher).Error

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

func TestSoftDeletePublisher(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		publisherID  uint
		mockBehavior func(mock sqlmock.Sqlmock, publisherID uint)
		wantErr      bool
	}{
		"Can soft delete a Publisher": {
			publisherID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, publisherID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `publishers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), publisherID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			publisherID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, publisherID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `publishers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Publisher"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.publisherID)

			err := db.Delete(&domain.Publisher{}, tc.publisherID).Error

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

func TestValidatePublisher(t *testing.T) {
	testCases := map[string]struct {
		publisher domain.Publisher
	}{
		"Can empty validations errors": {
			publisher: domain.Publisher{
				Name:   "Game Science",
				Acting: true,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.publisher.ValidatePublisher()
			assert.NoError(t, err)
		})
	}
}

func TestCreatePublisherWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		publisher domain.Publisher
		wantErr   string
	}{
		"Missing required fields": {
			publisher: domain.Publisher{},
			wantErr: `
				Name is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.publisher.ValidatePublisher()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
