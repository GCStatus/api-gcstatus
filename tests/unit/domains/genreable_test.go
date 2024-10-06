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

func TestCreateGenreable(t *testing.T) {
	testCases := map[string]struct {
		genreable    domain.Genreable
		mockBehavior func(mock sqlmock.Sqlmock, genreable domain.Genreable)
		expectError  bool
	}{
		"Success": {
			genreable: domain.Genreable{
				GenreableID:   1,
				GenreableType: "games",
				GenreID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, genreable domain.Genreable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `genreables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						genreable.GenreableID,
						genreable.GenreableType,
						genreable.GenreID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			genreable: domain.Genreable{
				GenreableID:   1,
				GenreableType: "games",
				GenreID:       1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, genreable domain.Genreable) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `genreables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						genreable.GenreableID,
						genreable.GenreableType,
						genreable.GenreID,
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

			tc.mockBehavior(mock, tc.genreable)

			err := db.Create(&tc.genreable).Error

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

func TestUpdateGenreable(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		genreable    domain.Genreable
		mockBehavior func(mock sqlmock.Sqlmock, genreable domain.Genreable)
		expectError  bool
	}{
		"Success": {
			genreable: domain.Genreable{
				ID:            1,
				GenreableID:   1,
				GenreableType: "games",
				GenreID:       1,
				CreatedAt:     fixedTime,
				UpdatedAt:     fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, genreable domain.Genreable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `genreables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						genreable.GenreableID,
						genreable.GenreableType,
						genreable.GenreID,
						genreable.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			genreable: domain.Genreable{
				ID:            1,
				GenreableID:   1,
				GenreableType: "games",
				GenreID:       1,
				CreatedAt:     fixedTime,
				UpdatedAt:     fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, genreable domain.Genreable) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `genreables`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						genreable.GenreableID,
						genreable.GenreableType,
						genreable.GenreID,
						genreable.ID,
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

			tc.mockBehavior(mock, tc.genreable)

			err := db.Save(&tc.genreable).Error

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

func TestSoftDeleteGenreable(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		genreableID  uint
		mockBehavior func(mock sqlmock.Sqlmock, genreableID uint)
		wantErr      bool
	}{
		"Can soft delete a Genreable": {
			genreableID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, genreableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `genreables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), genreableID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			genreableID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, genreableID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `genreables` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Genreable"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.genreableID)

			err := db.Delete(&domain.Genreable{}, tc.genreableID).Error

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
