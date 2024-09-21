package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"testing"
	"time"
)

func TestTransformProfile(t *testing.T) {
	staticTime := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		input    domain.Profile
		expected resources.ProfileResource
	}{
		"normal profile": {
			input: domain.Profile{
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
				CreatedAt: staticTime,
				UpdatedAt: staticTime,
				UserID:    1,
			},
			expected: resources.ProfileResource{
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
				CreatedAt: staticTime.Format("2006-01-02T15:04:05"),
				UpdatedAt: staticTime.Format("2006-01-02T15:04:05"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			profileResource := resources.TransformProfile(test.input)

			if profileResource.ID != test.expected.ID {
				t.Errorf("Expected ID %d, got %d", test.expected.ID, profileResource.ID)
			}
			if profileResource.Share != test.expected.Share {
				t.Errorf("Expected Share %v, got %v", test.expected.Share, profileResource.Share)
			}
			if profileResource.Photo != test.expected.Photo {
				t.Errorf("Expected Photo %s, got %s", test.expected.Photo, profileResource.Photo)
			}
			if profileResource.Phone != test.expected.Phone {
				t.Errorf("Expected Phone %s, got %s", test.expected.Phone, profileResource.Phone)
			}
			if profileResource.Facebook != test.expected.Facebook {
				t.Errorf("Expected Facebook %s, got %s", test.expected.Facebook, profileResource.Facebook)
			}
			if profileResource.Instagram != test.expected.Instagram {
				t.Errorf("Expected Instagram %s, got %s", test.expected.Instagram, profileResource.Instagram)
			}
			if profileResource.Twitter != test.expected.Twitter {
				t.Errorf("Expected Twitter %s, got %s", test.expected.Twitter, profileResource.Twitter)
			}
			if profileResource.Youtube != test.expected.Youtube {
				t.Errorf("Expected Youtube %s, got %s", test.expected.Youtube, profileResource.Youtube)
			}
			if profileResource.Twitch != test.expected.Twitch {
				t.Errorf("Expected Twitch %s, got %s", test.expected.Twitch, profileResource.Twitch)
			}
			if profileResource.Github != test.expected.Github {
				t.Errorf("Expected Github %s, got %s", test.expected.Github, profileResource.Github)
			}
			if profileResource.CreatedAt != test.expected.CreatedAt {
				t.Errorf("Expected CreatedAt %s, got %s", test.expected.CreatedAt, profileResource.CreatedAt)
			}
			if profileResource.UpdatedAt != test.expected.UpdatedAt {
				t.Errorf("Expected UpdatedAt %s, got %s", test.expected.UpdatedAt, profileResource.UpdatedAt)
			}
		})
	}
}
