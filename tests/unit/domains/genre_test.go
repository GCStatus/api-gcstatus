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

func TestCreateGenre(t *testing.T) {
	testCases := map[string]struct {
		genre        domain.Genre
		mockBehavior func(mock sqlmock.Sqlmock, genre domain.Genre)
		expectError  bool
	}{
		"Success": {
			genre: domain.Genre{
				Name: "Genre 1",
				Slug: "genre-1",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, genre domain.Genre) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `genres`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						genre.Name,
						genre.Slug,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			genre: domain.Genre{
				Name: "Failure",
				Slug: "failure",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, genre domain.Genre) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `genres`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						genre.Name,
						genre.Slug,
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

			tc.mockBehavior(mock, tc.genre)

			err := db.Create(&tc.genre).Error

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

func TestUpdateGenre(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		genre        domain.Genre
		mockBehavior func(mock sqlmock.Sqlmock, genre domain.Genre)
		expectError  bool
	}{
		"Success": {
			genre: domain.Genre{
				ID:        1,
				Name:      "Genre 1",
				Slug:      "genre-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, genre domain.Genre) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `genres`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						genre.Name,
						genre.Slug,
						genre.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			genre: domain.Genre{
				ID:        1,
				Name:      "Genre 1",
				Slug:      "genre-1",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, genre domain.Genre) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `genres`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						genre.Name,
						genre.Slug,
						genre.ID,
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

			tc.mockBehavior(mock, tc.genre)

			err := db.Save(&tc.genre).Error

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

func TestSoftDeleteGenre(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		genreID      uint
		mockBehavior func(mock sqlmock.Sqlmock, genreID uint)
		wantErr      bool
	}{
		"Can soft delete a Genre": {
			genreID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, genreID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `genres` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), genreID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			genreID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, genreID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `genres` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Genre"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.genreID)

			err := db.Delete(&domain.Genre{}, tc.genreID).Error

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

func TestGetGenreByID(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		genreID   uint
		mockFunc  func()
		wantGenre domain.Genre
		wantError bool
	}{
		"Valid Genre fetch": {
			genreID: 1,
			wantGenre: domain.Genre{
				ID:   1,
				Name: "Genre 1",
			},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Genre 1")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `genres` WHERE `genres`.`id` = ? AND `genres`.`deleted_at` IS NULL ORDER BY `genres`.`id` LIMIT ?")).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			wantError: false,
		},
		"Genre not found": {
			genreID:   2,
			wantGenre: domain.Genre{},
			wantError: true,
			mockFunc: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `genres` WHERE `genres`.`id` = ? AND `genres`.`deleted_at` IS NULL ORDER BY `genres`.`id` LIMIT ?")).
					WithArgs(2, 1).WillReturnError(fmt.Errorf("record not found"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			var genre domain.Genre
			err := db.First(&genre, tc.genreID).Error

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantGenre, genre)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidateGenreValidData(t *testing.T) {
	testCases := map[string]struct {
		genre domain.Genre
	}{
		"Can empty validations errors": {
			genre: domain.Genre{
				Name: "Genre 1",
				Slug: "genre-1",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.genre.ValidateGenre()
			assert.NoError(t, err)
		})
	}
}

func TestCreateGenreWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		genre   domain.Genre
		wantErr string
	}{
		"Missing required fields": {
			genre:   domain.Genre{},
			wantErr: "Name is a required field, Slug is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.genre.ValidateGenre()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
