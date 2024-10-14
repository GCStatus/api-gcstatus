package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"gcstatus/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateDLCLanguage(t *testing.T) {
	testCases := map[string]struct {
		DLCLanguage  domain.DLCLanguage
		mockBehavior func(mock sqlmock.Sqlmock, DLCLanguage domain.DLCLanguage)
		expectError  bool
	}{
		"Success": {
			DLCLanguage: domain.DLCLanguage{
				Menu:       false,
				Dubs:       true,
				Subtitles:  false,
				LanguageID: 1,
				DLCID:      1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCLanguage domain.DLCLanguage) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dlc_languages`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCLanguage.Menu,
						DLCLanguage.Dubs,
						DLCLanguage.Subtitles,
						DLCLanguage.LanguageID,
						DLCLanguage.DLCID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			DLCLanguage: domain.DLCLanguage{
				Menu:       false,
				Dubs:       true,
				Subtitles:  false,
				LanguageID: 1,
				DLCID:      1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCLanguage domain.DLCLanguage) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dlc_languages`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCLanguage.Menu,
						DLCLanguage.Dubs,
						DLCLanguage.Subtitles,
						DLCLanguage.LanguageID,
						DLCLanguage.DLCID,
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

			tc.mockBehavior(mock, tc.DLCLanguage)

			err := db.Create(&tc.DLCLanguage).Error

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

func TestUpdateDLCLanguage(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		DLCLanguage  domain.DLCLanguage
		mockBehavior func(mock sqlmock.Sqlmock, DLCLanguage domain.DLCLanguage)
		expectError  bool
	}{
		"Success": {
			DLCLanguage: domain.DLCLanguage{
				ID:         1,
				Menu:       false,
				Dubs:       true,
				Subtitles:  false,
				LanguageID: 1,
				DLCID:      1,
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCLanguage domain.DLCLanguage) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_languages`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCLanguage.Menu,
						DLCLanguage.Dubs,
						DLCLanguage.Subtitles,
						DLCLanguage.LanguageID,
						DLCLanguage.DLCID,
						DLCLanguage.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			DLCLanguage: domain.DLCLanguage{
				ID:         1,
				Menu:       false,
				Dubs:       true,
				Subtitles:  false,
				LanguageID: 1,
				DLCID:      1,
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCLanguage domain.DLCLanguage) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_languages`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCLanguage.Menu,
						DLCLanguage.Dubs,
						DLCLanguage.Subtitles,
						DLCLanguage.LanguageID,
						DLCLanguage.DLCID,
						DLCLanguage.ID,
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

			tc.mockBehavior(mock, tc.DLCLanguage)

			err := db.Save(&tc.DLCLanguage).Error

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

func TestSoftDeleteDLCLanguage(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		DLCLanguageID uint
		mockBehavior  func(mock sqlmock.Sqlmock, DLCLanguageID uint)
		wantErr       bool
	}{
		"Can soft delete a DLCLanguage": {
			DLCLanguageID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, DLCLanguageID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_languages` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), DLCLanguageID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			DLCLanguageID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, DLCLanguageID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_languages` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete DLCLanguage"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.DLCLanguageID)

			err := db.Delete(&domain.DLCLanguage{}, tc.DLCLanguageID).Error

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

func TestValidateDLCLanguageLanguageValidData(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		DLCLanguage domain.DLCLanguage
	}{
		"Can empty validations errors": {
			DLCLanguage: domain.DLCLanguage{
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
				DLC: domain.DLC{
					Name:             "Game Science",
					About:            "About DLC",
					Description:      "DLC Description",
					ShortDescription: "Short DLC Description",
					Cover:            "https://google.com",
					ReleaseDate:      fixedTime,
					Game: domain.Game{
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
						Views: []domain.Viewable{
							{
								UserID:       10,
								ViewableID:   1,
								ViewableType: "games",
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.DLCLanguage.ValidateDLCLanguage()
			assert.NoError(t, err)
		})
	}
}

func TestCreateDLCLanguageWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		DLCLanguage domain.DLCLanguage
		wantErr     string
	}{
		"Missing required fields": {
			DLCLanguage: domain.DLCLanguage{},
			wantErr: `
				Name is a required field,
				ISO is a required field,
				Name is a required field,
				Cover is a required field,
				About is a required field,
				Description is a required field,
				ShortDescription is a required field,
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
			err := tc.DLCLanguage.ValidateDLCLanguage()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
