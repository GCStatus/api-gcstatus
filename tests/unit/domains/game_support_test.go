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

func TestCreateGameSupport(t *testing.T) {
	testCases := map[string]struct {
		gameSupport  domain.GameSupport
		mockBehavior func(mock sqlmock.Sqlmock, gameSupport domain.GameSupport)
		expectError  bool
	}{
		"Success": {
			gameSupport: domain.GameSupport{
				URL:     utils.StringPtr("https://google.com"),
				Email:   utils.StringPtr("fake@example.com"),
				Contact: utils.StringPtr("fakecontact"),
				GameID:  1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameSupport domain.GameSupport) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `game_supports`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameSupport.URL,
						gameSupport.Email,
						gameSupport.Contact,
						gameSupport.GameID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			gameSupport: domain.GameSupport{
				URL:     utils.StringPtr("https://google.com"),
				Email:   utils.StringPtr("fake@example.com"),
				Contact: utils.StringPtr("fakecontact"),
				GameID:  1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameSupport domain.GameSupport) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `game_supports`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameSupport.URL,
						gameSupport.Email,
						gameSupport.Contact,
						gameSupport.GameID,
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

			tc.mockBehavior(mock, tc.gameSupport)

			err := db.Create(&tc.gameSupport).Error

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

func TestUpdateGameSupport(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		gameSupport  domain.GameSupport
		mockBehavior func(mock sqlmock.Sqlmock, gameSupport domain.GameSupport)
		expectError  bool
	}{
		"Success": {
			gameSupport: domain.GameSupport{
				ID:        1,
				URL:       utils.StringPtr("https://google.com"),
				Email:     utils.StringPtr("fake@example.com"),
				Contact:   utils.StringPtr("fakecontact"),
				GameID:    1,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameSupport domain.GameSupport) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_supports`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameSupport.URL,
						gameSupport.Email,
						gameSupport.Contact,
						gameSupport.GameID,
						gameSupport.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			gameSupport: domain.GameSupport{
				ID:        1,
				URL:       utils.StringPtr("https://google.com"),
				Email:     utils.StringPtr("fake@example.com"),
				Contact:   utils.StringPtr("fakecontact"),
				GameID:    1,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, gameSupport domain.GameSupport) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_supports`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						gameSupport.URL,
						gameSupport.Email,
						gameSupport.Contact,
						gameSupport.GameID,
						gameSupport.ID,
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

			tc.mockBehavior(mock, tc.gameSupport)

			err := db.Save(&tc.gameSupport).Error

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

func TestSoftDeleteGameSupport(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		gameSupportID uint
		mockBehavior  func(mock sqlmock.Sqlmock, gameSupportID uint)
		wantErr       bool
	}{
		"Can soft delete a GameSupport": {
			gameSupportID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, gameSupportID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_supports` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), gameSupportID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			gameSupportID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, GameSupport uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `game_supports` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete GameSupport"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.gameSupportID)

			err := db.Delete(&domain.GameSupport{}, tc.gameSupportID).Error

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

func TestValidateGameSupport(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		gameSupport domain.GameSupport
	}{
		"Can empty validations errors": {
			gameSupport: domain.GameSupport{
				URL:     utils.StringPtr("https://google.com"),
				Email:   utils.StringPtr("fake@example.com"),
				Contact: utils.StringPtr("fakecontact"),
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
			err := tc.gameSupport.ValidateGameSupport()
			assert.NoError(t, err)
		})
	}
}

func TestCreateGameSupportWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		gameSupport domain.GameSupport
		wantErr     string
	}{
		"Missing required fields": {
			gameSupport: domain.GameSupport{},
			wantErr: `
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
			err := tc.gameSupport.ValidateGameSupport()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
