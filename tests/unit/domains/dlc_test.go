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

func TestCreateDLC(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		DLC          domain.DLC
		mockBehavior func(mock sqlmock.Sqlmock, DLC domain.DLC)
		expectError  bool
	}{
		"Success": {
			DLC: domain.DLC{
				Name:             "Game Science",
				About:            "About DLC",
				Description:      "DLC Description",
				Free:             false,
				ShortDescription: "Short DLC Description",
				Cover:            "https://google.com",
				ReleaseDate:      fixedTime,
				GameID:           1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLC domain.DLC) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dlcs`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLC.Name,
						DLC.Cover,
						DLC.About,
						DLC.Description,
						DLC.Free,
						DLC.ShortDescription,
						DLC.Legal,
						DLC.ReleaseDate,
						DLC.GameID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			DLC: domain.DLC{
				Name:             "Game Science",
				About:            "About DLC",
				Description:      "DLC Description",
				ShortDescription: "Short DLC Description",
				Free:             false,
				Cover:            "https://google.com",
				ReleaseDate:      fixedTime,
				GameID:           1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLC domain.DLC) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dlcs`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLC.Name,
						DLC.Cover,
						DLC.About,
						DLC.Description,
						DLC.Free,
						DLC.ShortDescription,
						DLC.Legal,
						DLC.ReleaseDate,
						DLC.GameID,
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

			tc.mockBehavior(mock, tc.DLC)

			err := db.Create(&tc.DLC).Error

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

func TestUpdateDLC(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		DLC          domain.DLC
		mockBehavior func(mock sqlmock.Sqlmock, DLC domain.DLC)
		expectError  bool
	}{
		"Success": {
			DLC: domain.DLC{
				ID:               1,
				Name:             "Game Science",
				About:            "About DLC",
				Description:      "DLC Description",
				Free:             false,
				ShortDescription: "Short DLC Description",
				Cover:            "https://google.com",
				ReleaseDate:      fixedTime,
				GameID:           1,
				CreatedAt:        fixedTime,
				UpdatedAt:        fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLC domain.DLC) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlcs`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLC.Name,
						DLC.Cover,
						DLC.About,
						DLC.Description,
						DLC.Free,
						DLC.ShortDescription,
						DLC.Legal,
						DLC.ReleaseDate,
						DLC.GameID,
						DLC.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			DLC: domain.DLC{
				ID:               1,
				Name:             "Game Science",
				About:            "About DLC",
				Free:             false,
				Description:      "DLC Description",
				ShortDescription: "Short DLC Description",
				Cover:            "https://google.com",
				ReleaseDate:      fixedTime,
				GameID:           1,
				CreatedAt:        fixedTime,
				UpdatedAt:        fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLC domain.DLC) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlcs`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLC.Name,
						DLC.Cover,
						DLC.About,
						DLC.Description,
						DLC.Free,
						DLC.ShortDescription,
						DLC.Legal,
						DLC.ReleaseDate,
						DLC.GameID,
						DLC.ID,
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

			tc.mockBehavior(mock, tc.DLC)

			err := db.Save(&tc.DLC).Error

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

func TestSoftDeleteDLC(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		DLCID        uint
		mockBehavior func(mock sqlmock.Sqlmock, DLCID uint)
		wantErr      bool
	}{
		"Can soft delete a DLC": {
			DLCID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, DLCID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlcs` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), DLCID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			DLCID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, DLCID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlcs` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete DLC"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.DLCID)

			err := db.Delete(&domain.DLC{}, tc.DLCID).Error

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

func TestValidateDLC(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		DLC domain.DLC
	}{
		"Can empty validations errors": {
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
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.DLC.ValidateDLC()
			assert.NoError(t, err)
		})
	}
}

func TestCreateDLCWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		DLC     domain.DLC
		wantErr string
	}{
		"Missing required fields": {
			DLC: domain.DLC{},
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
				ReleaseDate is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.DLC.ValidateDLC()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
