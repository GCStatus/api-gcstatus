package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateCriticable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		criticable   domain.Criticable
		mockBehavior func(mock sqlmock.Sqlmock, criticable domain.Criticable)
		expectError  bool
	}{
		"Success": {
			criticable: domain.Criticable{
				Rate:           decimal.NewFromUint64(5),
				URL:            "https://google.com",
				PostedAt:       fixedTime,
				CriticID:       1,
				CriticableID:   1,
				CriticableType: "games",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, criticable domain.Criticable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `criticables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						criticable.Rate,
						criticable.URL,
						criticable.PostedAt,
						criticable.CriticableID,
						criticable.CriticableType,
						criticable.CriticID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			criticable: domain.Criticable{
				Rate:           decimal.NewFromUint64(5),
				URL:            "https://google.com",
				CriticID:       1,
				CriticableID:   1,
				CriticableType: "games",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, criticable domain.Criticable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `criticables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						criticable.Rate,
						criticable.URL,
						criticable.PostedAt,
						criticable.CriticableID,
						criticable.CriticableType,
						criticable.CriticID,
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

			tc.mockBehavior(mock, tc.criticable)

			err := db.Create(&tc.criticable).Error

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

func TestUpdateCriticable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		criticable   domain.Criticable
		mockBehavior func(mock sqlmock.Sqlmock, criticable domain.Criticable)
		expectError  bool
	}{
		"Success": {
			criticable: domain.Criticable{
				ID:             1,
				Rate:           decimal.NewFromUint64(5),
				URL:            "https://google.com",
				CriticID:       1,
				CriticableID:   1,
				CriticableType: "games",
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, criticable domain.Criticable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `criticables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						criticable.Rate,
						criticable.URL,
						criticable.PostedAt,
						criticable.CriticableID,
						criticable.CriticableType,
						criticable.CriticID,
						criticable.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			criticable: domain.Criticable{
				ID:             1,
				Rate:           decimal.NewFromUint64(5),
				URL:            "https://google.com",
				CriticID:       1,
				CriticableID:   1,
				CriticableType: "games",
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, criticable domain.Criticable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `criticables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						criticable.Rate,
						criticable.URL,
						criticable.PostedAt,
						criticable.CriticableID,
						criticable.CriticableType,
						criticable.CriticID,
						criticable.ID,
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

			tc.mockBehavior(mock, tc.criticable)

			err := db.Save(&tc.criticable).Error

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

func TestSoftDeleteCriticable(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		criticableID uint
		mockBehavior func(mock sqlmock.Sqlmock, criticableID uint)
		wantErr      bool
	}{
		"Can soft delete a Criticable": {
			criticableID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, criticableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `criticables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), criticableID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			criticableID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, criticableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `criticables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Criticable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.criticableID)

			err := db.Delete(&domain.Criticable{}, tc.criticableID).Error

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

func TestValidateCriticableValidData(t *testing.T) {
	testCases := map[string]struct {
		criticable domain.Criticable
	}{
		"Can empty validations errors": {
			criticable: domain.Criticable{
				Rate:           decimal.NewFromUint64(5),
				URL:            "https://google.com",
				CriticableID:   1,
				CriticableType: "games",
				Critic: domain.Critic{
					Name: "Test",
					URL:  "https://google.com",
					Logo: "https://placehold.co/600x400/EEE/31343C",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.criticable.ValidateCriticable()
			assert.NoError(t, err)
		})
	}
}

func TestCreateCriticableWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		criticable domain.Criticable
		wantErr    string
	}{
		"Missing required fields": {
			criticable: domain.Criticable{},
			wantErr:    "Name is a required field, URL is a required field, Logo is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.criticable.ValidateCriticable()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
