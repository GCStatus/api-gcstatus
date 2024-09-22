package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	"gcstatus/tests"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateProfile(t *testing.T) {
	testCases := map[string]struct {
		profile      domain.Profile
		mockBehavior func(mock sqlmock.Sqlmock, profile domain.Profile)
		expectError  bool
	}{
		"Success": {
			profile: domain.Profile{
				Share:     true,
				Photo:     "https://placehold.co/600x400/EEE/31343C",
				Phone:     "5511928342813",
				Facebook:  "https://facebook.com/any",
				Instagram: "https://instagram.com/any",
				Twitter:   "https://twitter.com/any",
				Youtube:   "https://youtube.com/any",
				Twitch:    "https://twitch.com/any",
				Github:    "https://github.com/any",
				UserID:    1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, profile domain.Profile) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `profiles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						profile.Share,
						profile.Photo,
						profile.Phone,
						profile.Facebook,
						profile.Instagram,
						profile.Twitter,
						profile.Youtube,
						profile.Twitch,
						profile.Github,
						profile.UserID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			profile: domain.Profile{
				Share:     true,
				Photo:     "https://placehold.co/600x400/EEE/31343C",
				Phone:     "5511928342813",
				Facebook:  "https://facebook.com/any",
				Instagram: "https://instagram.com/any",
				Twitter:   "https://twitter.com/any",
				Youtube:   "https://youtube.com/any",
				Twitch:    "https://twitch.com/any",
				Github:    "https://github.com/any",
				UserID:    1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, profile domain.Profile) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `profiles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						profile.Share,
						profile.Photo,
						profile.Phone,
						profile.Facebook,
						profile.Instagram,
						profile.Twitter,
						profile.Youtube,
						profile.Twitch,
						profile.Github,
						profile.UserID,
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

			tc.mockBehavior(mock, tc.profile)

			err := db.Create(&tc.profile).Error

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

func TestUpdateProfile(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		profile      domain.Profile
		mockBehavior func(mock sqlmock.Sqlmock, profile domain.Profile)
		expectError  bool
	}{
		"Success": {
			profile: domain.Profile{
				ID:        1,
				Share:     true,
				Photo:     "https://placehold.co/600x400/EEE/31343C",
				Phone:     "5511928342813",
				Facebook:  "https://facebook.com/any",
				Instagram: "https://instagram.com/any",
				Twitter:   "https://twitter.com/any",
				Youtube:   "https://youtube.com/any",
				Twitch:    "https://twitch.com/any",
				Github:    "https://github.com/any",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				UserID:    1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, profile domain.Profile) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `profiles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						profile.Share,
						profile.Photo,
						profile.Phone,
						profile.Facebook,
						profile.Instagram,
						profile.Twitter,
						profile.Youtube,
						profile.Twitch,
						profile.Github,
						profile.UserID,
						profile.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			profile: domain.Profile{
				ID:        1,
				Share:     true,
				Photo:     "https://placehold.co/600x400/EEE/31343C",
				Phone:     "5511928342813",
				Facebook:  "https://facebook.com/any",
				Instagram: "https://instagram.com/any",
				Twitter:   "https://twitter.com/any",
				Youtube:   "https://youtube.com/any",
				Twitch:    "https://twitch.com/any",
				Github:    "https://github.com/any",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				UserID:    1,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, profile domain.Profile) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `profiles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						profile.Share,
						profile.Photo,
						profile.Phone,
						profile.Facebook,
						profile.Instagram,
						profile.Twitter,
						profile.Youtube,
						profile.Twitch,
						profile.Github,
						profile.UserID,
						profile.ID,
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

			tc.mockBehavior(mock, tc.profile)

			err := db.Save(&tc.profile).Error

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

func TestSoftDeleteProfile(t *testing.T) {
	db, mock := tests.Setup(t)

	testCases := map[string]struct {
		profileID    uint
		mockBehavior func(mock sqlmock.Sqlmock, profileID uint)
		wantErr      bool
	}{
		"Can soft delete a level": {
			profileID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, profileID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `profiles` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), profileID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			profileID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, profileID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `profiles` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete profile"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.profileID)

			err := db.Delete(&domain.Profile{}, tc.profileID).Error

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

func TestValidateProfileValidData(t *testing.T) {
	testCases := map[string]struct {
		profile domain.Profile
	}{
		"Can empty validations errors": {
			profile: domain.Profile{
				ID:        1,
				Share:     true,
				Photo:     "https://placehold.co/600x400/EEE/31343C",
				Phone:     "5511928342813",
				Facebook:  "https://facebook.com/any",
				Instagram: "https://instagram.com/any",
				Twitter:   "https://twitter.com/any",
				Youtube:   "https://youtube.com/any",
				Twitch:    "https://twitch.com/any",
				Github:    "https://github.com/any",
				UserID:    1,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.profile.ValidateProfile()
			assert.NoError(t, err)
		})
	}
}

func TestCreateProfileWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		profile domain.Profile
		wantErr string
	}{
		"Missing required fields": {
			profile: domain.Profile{},
			wantErr: "Share is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.profile.ValidateProfile()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
