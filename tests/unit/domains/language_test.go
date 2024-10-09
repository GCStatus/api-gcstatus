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

func TestCreateLanguage(t *testing.T) {
	testCases := map[string]struct {
		language     domain.Language
		mockBehavior func(mock sqlmock.Sqlmock, language domain.Language)
		expectError  bool
	}{
		"Success": {
			language: domain.Language{
				Name: "Portuguese",
				ISO:  "pt_BR",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, language domain.Language) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `languages`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						language.Name,
						language.ISO,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			language: domain.Language{
				Name: "Portuguese",
				ISO:  "pt_BR",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, language domain.Language) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `languages`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						language.Name,
						language.ISO,
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

			tc.mockBehavior(mock, tc.language)

			err := db.Create(&tc.language).Error

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

func TestUpdateLanguage(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		language     domain.Language
		mockBehavior func(mock sqlmock.Sqlmock, language domain.Language)
		expectError  bool
	}{
		"Success": {
			language: domain.Language{
				ID:        1,
				Name:      "Portuguese",
				ISO:       "pt_BR",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, language domain.Language) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `languages`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						language.Name,
						language.ISO,
						language.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			language: domain.Language{
				ID:        1,
				Name:      "Portuguese",
				ISO:       "pt_BR",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, language domain.Language) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `languages`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						language.Name,
						language.ISO,
						language.ID,
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

			tc.mockBehavior(mock, tc.language)

			err := db.Save(&tc.language).Error

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

func TestSoftDeleteLanguage(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		languageID   uint
		mockBehavior func(mock sqlmock.Sqlmock, languageID uint)
		wantErr      bool
	}{
		"Can soft delete a Language": {
			languageID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, languageID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `languages` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), languageID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			languageID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, languageID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `languages` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Language"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.languageID)

			err := db.Delete(&domain.Language{}, tc.languageID).Error

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

func TestGetLanguageByID(t *testing.T) {
	fixedTime := time.Now()
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		languageID   uint
		mockFunc     func()
		wantLanguage domain.Language
		wantError    bool
	}{
		"Valid Language fetch": {
			languageID: 1,
			wantLanguage: domain.Language{
				ID:        1,
				Name:      "Portuguese",
				ISO:       "pt_BR",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "iso", "created_at", "updated_at"}).
					AddRow(1, "Portuguese", "pt_BR", fixedTime, fixedTime)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `languages` WHERE `languages`.`id` = ? AND `languages`.`deleted_at` IS NULL ORDER BY `languages`.`id` LIMIT ?")).
					WithArgs(1, 1).WillReturnRows(rows)
			},
			wantError: false,
		},
		"Language not found": {
			languageID:   2,
			wantLanguage: domain.Language{},
			wantError:    true,
			mockFunc: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `languages` WHERE `languages`.`id` = ? AND `languages`.`deleted_at` IS NULL ORDER BY `languages`.`id` LIMIT ?")).
					WithArgs(2, 1).WillReturnError(fmt.Errorf("record not found"))
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockFunc()

			var Language domain.Language
			err := db.First(&Language, tc.languageID).Error

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantLanguage, Language)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestValidateLanguageValidData(t *testing.T) {
	testCases := map[string]struct {
		language domain.Language
	}{
		"Can empty validations errors": {
			language: domain.Language{
				Name: "Portuguese",
				ISO:  "pt_BR",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.language.ValidateLanguage()
			assert.NoError(t, err)
		})
	}
}

func TestCreateLanguageWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		language domain.Language
		wantErr  string
	}{
		"Missing required fields": {
			language: domain.Language{},
			wantErr:  "Name is a required field, ISO is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.language.ValidateLanguage()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
