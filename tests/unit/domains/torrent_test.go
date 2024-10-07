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

func TestCreateTorrent(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		torrent      domain.Torrent
		mockBehavior func(mock sqlmock.Sqlmock, torrent domain.Torrent)
		expectError  bool
	}{
		"Success": {
			torrent: domain.Torrent{
				URL:               "https:google.com",
				PostedAt:          fixedTime,
				TorrentProviderID: 1,
				GameID:            1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, torrent domain.Torrent) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `torrents`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						torrent.URL,
						torrent.PostedAt,
						torrent.TorrentProviderID,
						torrent.GameID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			torrent: domain.Torrent{
				URL:               "https:google.com",
				PostedAt:          fixedTime,
				TorrentProviderID: 1,
				GameID:            1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, torrent domain.Torrent) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `torrents`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						torrent.URL,
						torrent.PostedAt,
						torrent.TorrentProviderID,
						torrent.GameID,
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

			tc.mockBehavior(mock, tc.torrent)

			err := db.Create(&tc.torrent).Error

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

func TestUpdateTorrent(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		torrent      domain.Torrent
		mockBehavior func(mock sqlmock.Sqlmock, torrent domain.Torrent)
		expectError  bool
	}{
		"Success": {
			torrent: domain.Torrent{
				ID:                1,
				URL:               "https:google.com",
				PostedAt:          fixedTime,
				TorrentProviderID: 1,
				GameID:            1,
				CreatedAt:         fixedTime,
				UpdatedAt:         fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, torrent domain.Torrent) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `torrents`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						torrent.URL,
						torrent.PostedAt,
						torrent.TorrentProviderID,
						torrent.GameID,
						torrent.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			torrent: domain.Torrent{
				ID:                1,
				URL:               "https:google.com",
				PostedAt:          fixedTime,
				TorrentProviderID: 1,
				GameID:            1,
				CreatedAt:         fixedTime,
				UpdatedAt:         fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, torrent domain.Torrent) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `torrents`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						torrent.URL,
						torrent.PostedAt,
						torrent.TorrentProviderID,
						torrent.GameID,
						torrent.ID,
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

			tc.mockBehavior(mock, tc.torrent)

			err := db.Save(&tc.torrent).Error

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

func TestSoftDeleteTorrent(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		torrentID    uint
		mockBehavior func(mock sqlmock.Sqlmock, torrentID uint)
		wantErr      bool
	}{
		"Can soft delete a Torrent": {
			torrentID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, torrentID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `torrents` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), torrentID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			torrentID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, torrentID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `torrents` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Torrent"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.torrentID)

			err := db.Delete(&domain.Torrent{}, tc.torrentID).Error

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

func TestValidateTorrent(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		torrent domain.Torrent
	}{
		"Can empty validations errors": {
			torrent: domain.Torrent{
				URL:      "https:google.com",
				PostedAt: fixedTime,
				TorrentProvider: domain.TorrentProvider{
					URL:  "http://google.com",
					Name: "Google",
				},
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
			err := tc.torrent.ValidateTorrent()
			assert.NoError(t, err)
		})
	}
}

func TestCreateTorrentWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		torrent domain.Torrent
		wantErr string
	}{
		"Missing required fields": {
			torrent: domain.Torrent{},
			wantErr: `
				URL is a required field,
				Name is a required field,
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
			err := tc.torrent.ValidateTorrent()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
