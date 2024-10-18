package tests

import (
	"fmt"
	"gcstatus/internal/domain"
	testutils "gcstatus/tests/utils"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateBanner(t *testing.T) {
	testCases := map[string]struct {
		banner       domain.Banner
		mockBehavior func(mock sqlmock.Sqlmock, banner domain.Banner)
		expectError  bool
	}{
		"Success": {
			banner: domain.Banner{
				Component:      "header-home",
				BannerableID:   1,
				BannerableType: "games",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, banner domain.Banner) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `banners`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						banner.Component,
						banner.BannerableID,
						banner.BannerableType,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Insert Error": {
			banner: domain.Banner{
				Component:      "header-home",
				BannerableID:   1,
				BannerableType: "games",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, banner domain.Banner) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `banners`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						banner.Component,
						banner.BannerableID,
						banner.BannerableType,
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

			tc.mockBehavior(mock, tc.banner)

			err := db.Create(&tc.banner).Error

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

func TestUpdateBanner(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		banner       domain.Banner
		mockBehavior func(mock sqlmock.Sqlmock, banner domain.Banner)
		expectError  bool
	}{
		"Success": {
			banner: domain.Banner{
				ID:             1,
				Component:      "header-home",
				BannerableID:   1,
				BannerableType: "games",
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, banner domain.Banner) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `banners`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						banner.Component,
						banner.BannerableID,
						banner.BannerableType,
						banner.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		"Failure - Update Error": {
			banner: domain.Banner{
				ID:             1,
				Component:      "header-home",
				BannerableID:   1,
				BannerableType: "games",
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			mockBehavior: func(mock sqlmock.Sqlmock, banner domain.Banner) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `banners`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						banner.Component,
						banner.BannerableID,
						banner.BannerableType,
						banner.ID,
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

			tc.mockBehavior(mock, tc.banner)

			err := db.Save(&tc.banner).Error

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

func TestSoftDeleteBanner(t *testing.T) {
	db, mock := testutils.Setup(t)

	testCases := map[string]struct {
		bannerID     uint
		mockBehavior func(mock sqlmock.Sqlmock, bannerID uint)
		wantErr      bool
	}{
		"Can soft delete a Banner": {
			bannerID: 1,
			mockBehavior: func(mock sqlmock.Sqlmock, bannerID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `banners` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), bannerID).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		"Soft delete fails": {
			bannerID: 2,
			mockBehavior: func(mock sqlmock.Sqlmock, bannerID uint) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `banners` SET `deleted_at`").WithArgs(sqlmock.AnyArg(), 2).
					WillReturnError(fmt.Errorf("failed to delete Banner"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.mockBehavior(mock, tc.bannerID)

			err := db.Delete(&domain.Banner{}, tc.bannerID).Error

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

func TestValidateBannerValidData(t *testing.T) {
	testCases := map[string]struct {
		banner domain.Banner
	}{
		"Can empty validations errors": {
			banner: domain.Banner{
				Component:      "header-home",
				BannerableID:   1,
				BannerableType: "games",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.banner.ValidateBanner()
			assert.NoError(t, err)
		})
	}
}

func TestCreateBannerWithMissingFields(t *testing.T) {
	testCases := map[string]struct {
		banner  domain.Banner
		wantErr string
	}{
		"Missing required fields": {
			banner:  domain.Banner{},
			wantErr: "Component is a required field",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.banner.ValidateBanner()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
