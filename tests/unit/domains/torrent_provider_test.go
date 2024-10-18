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

func TestCreateTorrentProvider(t *testing.T) {
	testCases := map[string]struct {
		torrentProvider domain.TorrentProvider
		mockBehavior    func(mock sqlmock.Sqlmock, torrentProvider domain.TorrentProvider)
		expectError     bool
	}{
		"Success": {
			torrentProvider: domain.TorrentProvider{
				URL:  "https:google.com",
				Name: "Google",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, torrentProvider domain.TorrentProvider) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `torrent_providers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						torrentProvider.URL,
						torrentProvider.Name,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			torrentProvider: domain.TorrentProvider{
				URL:  "https:google.com",
				Name: "Google",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, torrentProvider domain.TorrentProvider) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `torrent_providers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						torrentProvider.URL,
						torrentProvider.Name,
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

			tc.mockBehavior(mock, tc.torrentProvider)

			err := db.Create(&tc.torrentProvider).Error

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

func TestUpdateTorrentProvider(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		torrentProvider domain.TorrentProvider
		mockBehavior    func(mock sqlmock.Sqlmock, torrentProvider domain.TorrentProvider)
		expectError     bool
	}{
		"Success": {
			torrentProvider: domain.TorrentProvider{
				ID:        1,
				URL:       "https:google.com",
				Name:      "Google",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, torrentProvider domain.TorrentProvider) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `torrent_providers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						torrentProvider.URL,
						torrentProvider.Name,
						torrentProvider.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			torrentProvider: domain.TorrentProvider{
				ID:        1,
				URL:       "https:google.com",
				Name:      "Google",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, torrentProvider domain.TorrentProvider) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `torrent_providers`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						torrentProvider.URL,
						torrentProvider.Name,
						torrentProvider.ID,
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

			tc.mockBehavior(mock, tc.torrentProvider)

			err := db.Save(&tc.torrentProvider).Error

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

func TestSoftDeleteTorrentProvider(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		torrentProviderID uint
		mockBehavior      func(mock sqlmock.Sqlmock, torrentProviderID uint)
		wantErr           bool
	}{
		"Can soft delete a TorrentProvider": {
			torrentProviderID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, torrentProviderID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `torrent_providers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), torrentProviderID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			torrentProviderID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, torrentProviderID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `torrent_providers` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete TorrentProvider"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.torrentProviderID)

			err := db.Delete(&domain.TorrentProvider{}, tc.torrentProviderID).Error

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

func TestValidateTorrentProvider(t *testing.T) {
	testCases := map[string]struct {
		torrentProvider domain.TorrentProvider
	}{
		"Can empty validations errors": {
			torrentProvider: domain.TorrentProvider{
				URL:  "https:google.com",
				Name: "Google",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.torrentProvider.ValidateTorrentProvider()
			assert.NoError(t, err)
		})
	}
}

func TestCreateTorrentProviderWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		torrentProvider domain.TorrentProvider
		wantErr         string
	}{
		"Missing required fields": {
			torrentProvider: domain.TorrentProvider{},
			wantErr: `
				URL is a required field,
				Name is a required field
			`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.torrentProvider.ValidateTorrentProvider()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), utils.NormalizeWhitespace(tc.wantErr))
		})
	}
}
