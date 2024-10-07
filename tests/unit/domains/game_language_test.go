package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
	"gcstatus/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateGameLanguage(t *testing.T) {
	testCases := map[string]struct {
		gameLanguage domain.GameLanguage
		mockBehavior func(mock sqlmock.Sqlmock, gameLanguage domain.GameLanguage)
		expectError  bool
	}{
		"Success": {
			gameLanguage: domain.GameLanguage{
				Menu:       false,
				Dubs:       true,
				Subtitles:  false,
				LanguageID: 1,
				GameID:     1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameLanguage domain.GameLanguage) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `game_languages`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameLanguage.Menu,
						gameLanguage.Dubs,
						gameLanguage.Subtitles,
						gameLanguage.LanguageID,
						gameLanguage.GameID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			gameLanguage: domain.GameLanguage{
				Menu:       false,
				Dubs:       true,
				Subtitles:  false,
				LanguageID: 1,
				GameID:     1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameLanguage domain.GameLanguage) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `game_languages`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameLanguage.Menu,
						gameLanguage.Dubs,
						gameLanguage.Subtitles,
						gameLanguage.LanguageID,
						gameLanguage.GameID,
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

			tc.mockBehavior(mock, tc.gameLanguage)

			err := db.Create(&tc.gameLanguage).Error

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

func TestUpdateGameLanguage(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		gameLanguage domain.GameLanguage
		mockBehavior func(mock sqlmock.Sqlmock, gameLanguage domain.GameLanguage)
		expectError  bool
	}{
		"Success": {
			gameLanguage: domain.GameLanguage{
				ID:         1,
				Menu:       false,
				Dubs:       true,
				Subtitles:  false,
				LanguageID: 1,
				GameID:     1,
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameLanguage domain.GameLanguage) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_languages`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameLanguage.Menu,
						gameLanguage.Dubs,
						gameLanguage.Subtitles,
						gameLanguage.LanguageID,
						gameLanguage.GameID,
						gameLanguage.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			gameLanguage: domain.GameLanguage{
				ID:         1,
				Menu:       false,
				Dubs:       true,
				Subtitles:  false,
				LanguageID: 1,
				GameID:     1,
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameLanguage domain.GameLanguage) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_languages`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameLanguage.Menu,
						gameLanguage.Dubs,
						gameLanguage.Subtitles,
						gameLanguage.LanguageID,
						gameLanguage.GameID,
						gameLanguage.ID,
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

			tc.mockBehavior(mock, tc.gameLanguage)

			err := db.Save(&tc.gameLanguage).Error

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

func TestSoftDeleteGameLanguage(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		gameLanguageID uint
		mockBehavior   func(mock sqlmock.Sqlmock, gameLanguageID uint)
		wantErr        bool
	}{
		"Can soft delete a GameLanguage": {
			gameLanguageID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, gameLanguageID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_languages` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), gameLanguageID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			gameLanguageID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, gameLanguageID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_languages` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete GameLanguage"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.gameLanguageID)

			err := db.Delete(&domain.GameLanguage{}, tc.gameLanguageID).Error

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

func TestValidateGameLanguageLanguageValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		gameLanguage domain.GameLanguage
	}{
		"Can empty validations errors": {
			gameLanguage: domain.GameLanguage{
				Menu:       false,
				Dubs:       true,
				Subtitles:  false,
				LanguageID: 1,
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
				Language: domain.Language{
					ID:        1,
					Name:      "Portuguese",
					ISO:       "pt_BR",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				Game: domain.Game{
					ID:               1,
					Slug:             "valid",
					Age:              18,
					Title:            "Game Test",
					Condition:        domain.CommomCondition,
					Cover:            "https://placehold.co/600x400/EEE/31343C",
					About:            "About game",
					Description:      "Description",
					ShortDescription: "Short description",
					Free:             false,
					ReleaseDate:      fixedTime,
					CreatedAt:        fixedTime,
					UpdatedAt:        fixedTime,
					View: domain.Viewable{
						Count:        10,
						ViewableID:   1,
						ViewableType: "games",
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.gameLanguage.ValidateGameLanguage()
			assert.NoError(t, err)
		})
	}
}

func TestCreateGameLanguageWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		gameLanguage domain.GameLanguage
		wantErr      string
	}{
		"Missing required fields": {
			gameLanguage: domain.GameLanguage{},
			wantErr: `
				Name is a required field,
				ISO is a required field,
				Age is a required field,
				Slug is a required field,
				Title is a required field,
				Condition is a required field,
				Cover is a required field,
				About is a required field,
				Description is a required field,
				ShortDescription is a required field,
				ReleaseDate is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.gameLanguage.ValidateGameLanguage()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
