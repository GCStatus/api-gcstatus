package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
	"gcstatus/pkg/utils"
	"testing"
)

type MockProfileRepository struct {
	profiles map[uint]*domain.Profile
}

func NewMockProfileRepository() *MockProfileRepository {
	return &MockProfileRepository{
		profiles: make(map[uint]*domain.Profile),
	}
}

func (m *MockProfileRepository) UpdatePicture(profileID uint, path string) error {
	profile, exists := m.profiles[profileID]
	if !exists {
		return errors.New("profile not found")
	}
	profile.Photo = path
	m.profiles[profileID] = profile
	return nil
}

func (m *MockProfileRepository) UpdateSocials(profileID uint, request ports.UpdateSocialsRequest) error {
	profile, exists := m.profiles[profileID]
	if !exists {
		return errors.New("profile not found")
	}

	if request.Share != nil {
		profile.Share = *request.Share
	}
	if request.Phone != nil {
		profile.Phone = *request.Phone
	}
	if request.Github != nil {
		profile.Github = *request.Github
	}
	if request.Twitch != nil {
		profile.Twitch = *request.Twitch
	}
	if request.Twitter != nil {
		profile.Twitter = *request.Twitter
	}
	if request.Youtube != nil {
		profile.Youtube = *request.Youtube
	}
	if request.Facebook != nil {
		profile.Facebook = *request.Facebook
	}
	if request.Instagram != nil {
		profile.Instagram = *request.Instagram
	}

	m.profiles[profileID] = profile
	return nil
}

func TestMockProfileRepository_UpdatePicture(t *testing.T) {
	mockRepo := NewMockProfileRepository()

	mockRepo.profiles[1] = &domain.Profile{
		ID:    1,
		Photo: "old-path.jpg",
	}

	testCases := map[string]struct {
		profileID     uint
		newPath       string
		expectedError bool
	}{
		"valid profile": {
			profileID:     1,
			newPath:       "new-path.jpg",
			expectedError: false,
		},
		"invalid profile": {
			profileID:     999,
			newPath:       "new-path.jpg",
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.UpdatePicture(tc.profileID, tc.newPath)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.profiles[tc.profileID].Photo != tc.newPath {
					t.Fatalf("expected photo path to be updated to %s, but got %s", tc.newPath, mockRepo.profiles[tc.profileID].Photo)
				}
			}
		})
	}
}

func TestMockProfileRepository_UpdateSocials(t *testing.T) {
	mockRepo := NewMockProfileRepository()

	mockRepo.profiles[1] = &domain.Profile{
		ID:        1,
		Share:     false,
		Phone:     "1234567890",
		Facebook:  "old-facebook",
		Instagram: "old-instagram",
	}

	testCases := map[string]struct {
		profileID     uint
		request       ports.UpdateSocialsRequest
		expectedError bool
	}{
		"valid update": {
			profileID: 1,
			request: ports.UpdateSocialsRequest{
				Share:     utils.BoolPtr(true),
				Phone:     utils.StringPtr("0987654321"),
				Facebook:  utils.StringPtr("new-facebook"),
				Instagram: utils.StringPtr("new-instagram"),
			},
			expectedError: false,
		},
		"non-existent profile": {
			profileID:     999,
			request:       ports.UpdateSocialsRequest{},
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.UpdateSocials(tc.profileID, tc.request)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				updatedProfile := mockRepo.profiles[tc.profileID]
				if updatedProfile.Share != *tc.request.Share {
					t.Fatalf("expected Share to be %v, got %v", *tc.request.Share, updatedProfile.Share)
				}
				if updatedProfile.Phone != *tc.request.Phone {
					t.Fatalf("expected Phone to be %s, got %s", *tc.request.Phone, updatedProfile.Phone)
				}
				if updatedProfile.Facebook != *tc.request.Facebook {
					t.Fatalf("expected Facebook to be %s, got %s", *tc.request.Facebook, updatedProfile.Facebook)
				}
				if updatedProfile.Instagram != *tc.request.Instagram {
					t.Fatalf("expected Instagram to be %s, got %s", *tc.request.Instagram, updatedProfile.Instagram)
				}
			}
		})
	}
}
