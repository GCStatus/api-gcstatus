package tests

import (
	"fmt"
	"regexp"
	"testing"

	"gcstatus/internal/adapters/db"
	"gcstatus/internal/ports"
	"gcstatus/pkg/utils"
	"gcstatus/tests"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUpdateSocials(t *testing.T) {
	gormDB, mock := tests.Setup(t)

	repo := db.NewProfileRepositoryMySQL(gormDB)

	tests := map[string]struct {
		profileID uint
		request   ports.UpdateSocialsRequest
		mock      func()
		expectErr bool
	}{
		"successful update": {
			profileID: 1,
			request: ports.UpdateSocialsRequest{
				Share:     utils.BoolPtr(true),
				Phone:     utils.StringPtr("123456789"),
				Github:    utils.StringPtr("githubUser"),
				Twitch:    nil,
				Twitter:   utils.StringPtr("twitterUser"),
				Youtube:   nil,
				Facebook:  utils.StringPtr("facebookUser"),
				Instagram: utils.StringPtr("instagramUser"),
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `profiles` SET `facebook`=?,`github`=?,`instagram`=?,`phone`=?,`share`=?,`twitch`=?,`twitter`=?,`youtube`=?,`updated_at`=? WHERE id = ? AND `profiles`.`deleted_at` IS NULL")).
					WithArgs("facebookUser", "githubUser", "instagramUser", "123456789", true, nil, "twitterUser", nil, sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		"failed update due to database error": {
			profileID: 2,
			request: ports.UpdateSocialsRequest{
				Share:     utils.BoolPtr(false),
				Phone:     nil,
				Github:    nil,
				Twitch:    nil,
				Twitter:   nil,
				Youtube:   nil,
				Facebook:  nil,
				Instagram: nil,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `profiles` SET `facebook`=?,`github`=?,`instagram`=?,`phone`=?,`share`=?,`twitch`=?,`twitter`=?,`youtube`=?,`updated_at`=? WHERE id = ? AND `profiles`.`deleted_at` IS NULL")).
					WithArgs(nil, nil, nil, nil, false, nil, nil, nil, sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.mock()
			err := repo.UpdateSocials(tt.profileID, tt.request)

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdatePicture(t *testing.T) {
	gormDB, mock := tests.Setup(t)

	repo := db.NewProfileRepositoryMySQL(gormDB)

	tests := map[string]struct {
		profileID uint
		path      string
		mock      func()
		expectErr bool
	}{
		"successful picture update": {
			profileID: 1,
			path:      "/path/to/photo.jpg",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `profiles`").
					WithArgs(
						"/path/to/photo.jpg",
						sqlmock.AnyArg(),
						1,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		"failed picture update due to database error": {
			profileID: 2,
			path:      "/path/to/photo.jpg",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `profiles`").
					WithArgs(
						"/path/to/photo.jpg",
						sqlmock.AnyArg(),
						2,
					).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.mock()
			err := repo.UpdatePicture(tt.profileID, tt.path)

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
