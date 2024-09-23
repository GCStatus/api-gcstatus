package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"testing"
	"time"
)

func TestTransformUser(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		inputUser domain.User
		expected  resources.UserResource
	}{
		"normal user": {
			inputUser: domain.User{
				ID:         1,
				Name:       "John Doe",
				Email:      "john@example.com",
				Nickname:   "Johnny",
				Blocked:    false,
				Birthdate:  fixedTime,
				Experience: 500,
				Coins:      1029,
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
				Profile:    domain.Profile{ID: 1, Photo: "https://google.com"},
				Level:      domain.Level{ID: 1, Level: 1, Experience: 0, Coins: 0},
			},
			expected: resources.UserResource{
				ID:         1,
				Name:       "John Doe",
				Email:      "john@example.com",
				Nickname:   "Johnny",
				Level:      1,
				Experience: 500,
				Coins:      1029,
				Birthdate:  fixedTime.Format("2006-01-02T15:04:05"),
				CreatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
				UpdatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
				Profile: &resources.ProfileResource{
					ID:    1,
					Photo: "https://google.com",
				},
			},
		},
		"empty user profile": {
			inputUser: domain.User{
				ID:         2,
				Name:       "Jane Smith",
				Email:      "jane@example.com",
				Nickname:   "Janey",
				Experience: 500,
				Coins:      1029,
				Blocked:    true,
				Birthdate:  fixedTime,
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
				Profile:    domain.Profile{ID: 0},
				Level:      domain.Level{ID: 1, Level: 1, Experience: 0, Coins: 0},
			},
			expected: resources.UserResource{
				ID:         2,
				Name:       "Jane Smith",
				Email:      "jane@example.com",
				Nickname:   "Janey",
				Experience: 500,
				Coins:      1029,
				Birthdate:  fixedTime.Format("2006-01-02T15:04:05"),
				CreatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
				UpdatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
				Profile:    nil,
				Level:      1,
			},
		},
		"empty user level": {
			inputUser: domain.User{
				ID:         2,
				Name:       "Jane Smith",
				Email:      "jane@example.com",
				Nickname:   "Janey",
				Blocked:    true,
				Experience: 500,
				Coins:      1029,
				Birthdate:  fixedTime,
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
				Profile:    domain.Profile{ID: 1, Photo: "https://google.com"},
				Level:      domain.Level{ID: 0},
			},
			expected: resources.UserResource{
				ID:         2,
				Name:       "Jane Smith",
				Email:      "jane@example.com",
				Nickname:   "Janey",
				Level:      0,
				Experience: 500,
				Coins:      1029,
				Birthdate:  fixedTime.Format("2006-01-02T15:04:05"),
				CreatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
				UpdatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
				Profile: &resources.ProfileResource{
					ID:    1,
					Photo: "https://google.com",
				},
			},
		},
		"missing name": {
			inputUser: domain.User{
				ID:         3,
				Name:       "",
				Email:      "no-name@example.com",
				Nickname:   "NoName",
				Blocked:    false,
				Experience: 500,
				Coins:      1029,
				Birthdate:  time.Date(1985, time.March, 15, 0, 0, 0, 0, time.UTC),
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
				Profile:    domain.Profile{ID: 2},
			},
			expected: resources.UserResource{
				ID:         3,
				Name:       "",
				Email:      "no-name@example.com",
				Experience: 500,
				Coins:      1029,
				Nickname:   "NoName",
				Birthdate:  time.Date(1985, time.March, 15, 0, 0, 0, 0, time.UTC).Format("2006-01-02T15:04:05"),
				CreatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
				UpdatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
				Profile:    &resources.ProfileResource{ID: 2},
			},
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			userResource := resources.TransformUser(test.inputUser)

			if userResource.ID != test.expected.ID {
				t.Errorf("Expected ID %d, got %d", test.expected.ID, userResource.ID)
			}
			if userResource.Name != test.expected.Name {
				t.Errorf("Expected Name %s, got %s", test.expected.Name, userResource.Name)
			}
			if userResource.Email != test.expected.Email {
				t.Errorf("Expected Email %s, got %s", test.expected.Email, userResource.Email)
			}
			if userResource.Nickname != test.expected.Nickname {
				t.Errorf("Expected Nickname %s, got %s", test.expected.Nickname, userResource.Nickname)
			}
			if userResource.Birthdate != test.expected.Birthdate {
				t.Errorf("Expected Birthdate %s, got %s", test.expected.Birthdate, userResource.Birthdate)
			}
			if userResource.CreatedAt != test.expected.CreatedAt {
				t.Errorf("Expected CreatedAt %s, got %s", test.expected.CreatedAt, userResource.CreatedAt)
			}
			if userResource.UpdatedAt != test.expected.UpdatedAt {
				t.Errorf("Expected UpdatedAt %s, got %s", test.expected.UpdatedAt, userResource.UpdatedAt)
			}
			if test.expected.Profile == nil && userResource.Profile != nil {
				t.Errorf("Expected Profile to be nil, got %+v", userResource.Profile)
			}
			if test.expected.Profile != nil && userResource.Profile != nil {
				if userResource.Profile.ID != test.expected.Profile.ID {
					t.Errorf("Expected Profile ID %d, got %d", test.expected.Profile.ID, userResource.Profile.ID)
				}
			}
		})
	}
}

func TestTransformUsers(t *testing.T) {
	fixedTime := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		inputUsers []domain.User
		expected   []resources.UserResource
	}{
		"two users": {
			inputUsers: []domain.User{
				{
					ID:         1,
					Name:       "John Doe",
					Email:      "john@example.com",
					Nickname:   "Johnny",
					Blocked:    false,
					Experience: 500,
					Coins:      1029,
					Birthdate:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
					CreatedAt:  fixedTime,
					UpdatedAt:  fixedTime,
				},
				{
					ID:         2,
					Name:       "Jane Smith",
					Email:      "jane@example.com",
					Nickname:   "Janey",
					Blocked:    true,
					Experience: 500,
					Coins:      1029,
					Birthdate:  time.Date(1985, time.March, 15, 0, 0, 0, 0, time.UTC),
					CreatedAt:  fixedTime,
					UpdatedAt:  fixedTime,
				},
			},
			expected: []resources.UserResource{
				{
					ID:         1,
					Name:       "John Doe",
					Email:      "john@example.com",
					Nickname:   "Johnny",
					Experience: 500,
					Coins:      1029,
					Birthdate:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02T15:04:05"),
					CreatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
					UpdatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
				},
				{
					ID:         2,
					Name:       "Jane Smith",
					Email:      "jane@example.com",
					Nickname:   "Janey",
					Experience: 500,
					Coins:      1029,
					Birthdate:  time.Date(1985, time.March, 15, 0, 0, 0, 0, time.UTC).Format("2006-01-02T15:04:05"),
					CreatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
					UpdatedAt:  fixedTime.Format("2006-01-02T15:04:05"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			userResources := resources.TransformUsers(test.inputUsers)

			if len(userResources) != len(test.expected) {
				t.Errorf("Expected %d resources, got %d", len(test.expected), len(userResources))
			}

			for i, user := range test.inputUsers {
				if userResources[i].ID != user.ID {
					t.Errorf("Expected ID %d, got %d", user.ID, userResources[i].ID)
				}
				if userResources[i].Name != user.Name {
					t.Errorf("Expected Name %s, got %s", user.Name, userResources[i].Name)
				}
				if userResources[i].Email != user.Email {
					t.Errorf("Expected Email %s, got %s", user.Email, userResources[i].Email)
				}
				if userResources[i].Nickname != user.Nickname {
					t.Errorf("Expected Nickname %s, got %s", user.Nickname, userResources[i].Nickname)
				}
				if userResources[i].Birthdate != user.Birthdate.Format("2006-01-02T15:04:05") {
					t.Errorf("Expected Birthdate %s, got %s", user.Birthdate, userResources[i].Birthdate)
				}
				if userResources[i].CreatedAt != user.CreatedAt.Format("2006-01-02T15:04:05") {
					t.Errorf("Expected CreatedAt %s, got %s", user.CreatedAt, userResources[i].CreatedAt)
				}
				if userResources[i].UpdatedAt != user.UpdatedAt.Format("2006-01-02T15:04:05") {
					t.Errorf("Expected UpdatedAt %s, got %s", user.UpdatedAt, userResources[i].UpdatedAt)
				}
				if userResources[i].Level != user.Level.Level {
					t.Errorf("Expected Level %v, got %v", user.Level, userResources[i].Level)
				}
			}
		})
	}
}