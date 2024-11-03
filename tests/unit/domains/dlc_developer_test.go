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

func TestCreateDLCDeveloper(t *testing.T) {
	testCases := map[string]struct {
		DLCDeveloper domain.DLCDeveloper
		mockBehavior func(mock sqlmock.Sqlmock, DLCDeveloper domain.DLCDeveloper)
		expectError  bool
	}{
		"Success": {
			DLCDeveloper: domain.DLCDeveloper{
				DLCID:       1,
				DeveloperID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCDeveloper domain.DLCDeveloper) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dlc_developers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCDeveloper.DLCID,
						DLCDeveloper.DeveloperID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			DLCDeveloper: domain.DLCDeveloper{
				DLCID:       1,
				DeveloperID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCDeveloper domain.DLCDeveloper) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dlc_developers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCDeveloper.DLCID,
						DLCDeveloper.DeveloperID,
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

			tc.mockBehavior(mock, tc.DLCDeveloper)

			err := db.Create(&tc.DLCDeveloper).Error

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

func TestUpdateDLCDeveloper(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		DLCDeveloper domain.DLCDeveloper
		mockBehavior func(mock sqlmock.Sqlmock, DLCDeveloper domain.DLCDeveloper)
		expectError  bool
	}{
		"Success": {
			DLCDeveloper: domain.DLCDeveloper{
				ID:          1,
				DLCID:       1,
				DeveloperID: 1,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCDeveloper domain.DLCDeveloper) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_developers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCDeveloper.DLCID,
						DLCDeveloper.DeveloperID,
						DLCDeveloper.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			DLCDeveloper: domain.DLCDeveloper{
				ID:          1,
				DLCID:       1,
				DeveloperID: 1,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCDeveloper domain.DLCDeveloper) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_developers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCDeveloper.DLCID,
						DLCDeveloper.DeveloperID,
						DLCDeveloper.ID,
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

			tc.mockBehavior(mock, tc.DLCDeveloper)

			err := db.Save(&tc.DLCDeveloper).Error

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

func TestSoftDeleteDLCDeveloper(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		DLCDeveloperID uint
		mockBehavior   func(mock sqlmock.Sqlmock, DLCDeveloperID uint)
		wantErr        bool
	}{
		"Can soft delete a DLCDeveloper": {
			DLCDeveloperID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, DLCDeveloperID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_developers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), DLCDeveloperID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			DLCDeveloperID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, DLCDeveloper uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_developers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete DLCDeveloper"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.DLCDeveloperID)

			err := db.Delete(&domain.DLCDeveloper{}, tc.DLCDeveloperID).Error

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

func TestValidateDLCDeveloper(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		DLCDeveloper domain.DLCDeveloper
	}{
		"Can empty validations errors": {
			DLCDeveloper: domain.DLCDeveloper{
				Developer: domain.Developer{
					Name:      "Game Science",
					Slug:      "game-science",
					Acting:    true,
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
			err := tc.DLCDeveloper.ValidateDLCDeveloper()
			assert.NoError(t, err)
		})
	}
}

func TestCreateDLCDeveloperWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		DLCDeveloper domain.DLCDeveloper
		wantErr      string
	}{
		"Missing required fields": {
			DLCDeveloper: domain.DLCDeveloper{},
			wantErr: `
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
				ReleaseDate is a required field,
				Name is a required field,
				Slug is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.DLCDeveloper.ValidateDLCDeveloper()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
