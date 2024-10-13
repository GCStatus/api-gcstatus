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

func TestCreateDLCPublisher(t *testing.T) {
	testCases := map[string]struct {
		DLCPublisher domain.DLCPublisher
		mockBehavior func(mock sqlmock.Sqlmock, DLCPublisher domain.DLCPublisher)
		expectError  bool
	}{
		"Success": {
			DLCPublisher: domain.DLCPublisher{
				DLCID:       1,
				PublisherID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCPublisher domain.DLCPublisher) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dlc_publishers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCPublisher.DLCID,
						DLCPublisher.PublisherID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			DLCPublisher: domain.DLCPublisher{
				DLCID:       1,
				PublisherID: 1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCPublisher domain.DLCPublisher) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dlc_publishers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCPublisher.DLCID,
						DLCPublisher.PublisherID,
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

			tc.mockBehavior(mock, tc.DLCPublisher)

			err := db.Create(&tc.DLCPublisher).Error

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

func TestUpdateDLCPublisher(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		DLCPublisher domain.DLCPublisher
		mockBehavior func(mock sqlmock.Sqlmock, DLCPublisher domain.DLCPublisher)
		expectError  bool
	}{
		"Success": {
			DLCPublisher: domain.DLCPublisher{
				ID:          1,
				DLCID:       1,
				PublisherID: 1,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCPublisher domain.DLCPublisher) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_publishers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCPublisher.DLCID,
						DLCPublisher.PublisherID,
						DLCPublisher.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			DLCPublisher: domain.DLCPublisher{
				ID:          1,
				DLCID:       1,
				PublisherID: 1,
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, DLCPublisher domain.DLCPublisher) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_publishers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						DLCPublisher.DLCID,
						DLCPublisher.PublisherID,
						DLCPublisher.ID,
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

			tc.mockBehavior(mock, tc.DLCPublisher)

			err := db.Save(&tc.DLCPublisher).Error

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

func TestSoftDeleteDLCPublisher(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		DLCPublisherID uint
		mockBehavior   func(mock sqlmock.Sqlmock, DLCPublisherID uint)
		wantErr        bool
	}{
		"Can soft delete a DLCPublisher": {
			DLCPublisherID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, DLCPublisherID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_publishers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), DLCPublisherID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			DLCPublisherID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, DLCPublisher uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dlc_publishers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete DLCPublisher"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.DLCPublisherID)

			err := db.Delete(&domain.DLCPublisher{}, tc.DLCPublisherID).Error

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

func TestValidateDLCPublisher(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		DLCPublisher domain.DLCPublisher
	}{
		"Can empty validations errors": {
			DLCPublisher: domain.DLCPublisher{
				Publisher: domain.Publisher{
					Name:      "Game Science",
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
			err := tc.DLCPublisher.ValidateDLCPublisher()
			assert.NoError(t, err)
		})
	}
}

func TestCreateDLCPublisherWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		DLCPublisher domain.DLCPublisher
		wantErr      string
	}{
		"Missing required fields": {
			DLCPublisher: domain.DLCPublisher{},
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
				Name is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.DLCPublisher.ValidateDLCPublisher()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
