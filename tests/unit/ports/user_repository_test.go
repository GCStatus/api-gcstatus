package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"testing"
)

type MockUserRepository struct {
	users map[uint]*domain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uint]*domain.User),
	}
}

func (m *MockUserRepository) CreateUser(user *domain.User) error {
	if user == nil {
		return errors.New("invalid user data")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetUserByID(id uint) (*domain.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) GetAllUsers() ([]domain.User, error) {
	var users []domain.User
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, nil
}

func (m *MockUserRepository) FindUserByEmailOrNickname(emailOrNickname string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == emailOrNickname || user.Nickname == emailOrNickname {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) UpdateUserPassword(userID uint, hashedPassword string) error {
	user, exists := m.users[userID]
	if !exists {
		return errors.New("user not found")
	}
	user.Password = hashedPassword
	return nil
}

func (m *MockUserRepository) CreateWithProfile(user *domain.User) error {
	if user == nil || user.Profile.ID == 0 {
		return errors.New("invalid user or profile data")
	}
	m.users[user.ID] = user
	return nil
}

func TestMockUserRepository_CreateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()

	testCases := map[string]struct {
		input         *domain.User
		expectedError bool
	}{
		"valid user": {
			input: &domain.User{
				ID:       1,
				Email:    "user@example.com",
				Nickname: "user1",
				Password: "hashedPassword",
			},
			expectedError: false,
		},
		"nil user": {
			input:         nil,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.CreateUser(tc.input)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.users[tc.input.ID] == nil {
					t.Fatalf("expected user to be created, but it wasn't")
				}
			}
		})
	}
}

func TestMockUserRepository_GetUserByID(t *testing.T) {
	mockRepo := NewMockUserRepository()
	err := mockRepo.CreateUser(&domain.User{
		ID:       1,
		Email:    "user@example.com",
		Nickname: "user1",
		Password: "hashedPassword",
	})

	if err != nil {
		t.Fatalf("failed to create the user: %s", err.Error())
	}

	testCases := map[string]struct {
		userID        uint
		expectedError bool
	}{
		"valid user ID": {
			userID:        1,
			expectedError: false,
		},
		"invalid user ID": {
			userID:        999,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			user, err := mockRepo.GetUserByID(tc.userID)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if user != nil {
					t.Fatalf("expected nil user, got %v", user)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if user == nil || user.ID != tc.userID {
					t.Fatalf("expected user ID %d, got %v", tc.userID, user)
				}
			}
		})
	}
}

func TestMockUserRepository_GetAllUsers(t *testing.T) {
	testCases := map[string]struct {
		expectedUserCount int
		mockCreateUser    func(repo *MockUserRepository)
	}{
		"multiple users": {
			expectedUserCount: 2,
			mockCreateUser: func(repo *MockUserRepository) {
				err := repo.CreateUser(&domain.User{
					ID:       1,
					Email:    "user1@example.com",
					Nickname: "user1",
					Password: "hashedPassword",
				})
				if err != nil {
					t.Fatalf("failed to create the user: %s", err.Error())
				}
				err = repo.CreateUser(&domain.User{
					ID:       2,
					Email:    "user2@example.com",
					Nickname: "user2",
					Password: "hashedPassword",
				})
				if err != nil {
					t.Fatalf("failed to create the user: %s", err.Error())
				}
			},
		},
		"no users": {
			expectedUserCount: 0,
			mockCreateUser:    func(repo *MockUserRepository) {},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockUserRepository()

			tc.mockCreateUser(mockRepo)

			users, err := mockRepo.GetAllUsers()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(users) != tc.expectedUserCount {
				t.Fatalf("expected %d users, got %d", tc.expectedUserCount, len(users))
			}
		})
	}
}

func TestMockUserRepository_FindUserByEmailOrNickname(t *testing.T) {
	mockRepo := NewMockUserRepository()
	err := mockRepo.CreateUser(&domain.User{
		ID:       1,
		Email:    "user@example.com",
		Nickname: "user1",
		Password: "hashedPassword",
	})

	if err != nil {
		t.Fatalf("failed to create the user: %s", err.Error())
	}

	testCases := map[string]struct {
		input         string
		expectedError bool
	}{
		"valid email": {
			input:         "user@example.com",
			expectedError: false,
		},
		"valid nickname": {
			input:         "user1",
			expectedError: false,
		},
		"invalid email or nickname": {
			input:         "invalid@example.com",
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			user, err := mockRepo.FindUserByEmailOrNickname(tc.input)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if user != nil {
					t.Fatalf("expected nil user, got %v", user)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if user == nil || (user.Email != tc.input && user.Nickname != tc.input) {
					t.Fatalf("expected user with email/nickname %s, got %v", tc.input, user)
				}
			}
		})
	}
}

func TestMockUserRepository_UpdateUserPassword(t *testing.T) {
	mockRepo := NewMockUserRepository()
	err := mockRepo.CreateUser(&domain.User{
		ID:       1,
		Email:    "user@example.com",
		Nickname: "user1",
		Password: "oldPassword",
	})

	if err != nil {
		t.Fatalf("failed to create the user: %s", err.Error())
	}

	testCases := map[string]struct {
		userID        uint
		newPassword   string
		expectedError bool
	}{
		"valid user ID": {
			userID:        1,
			newPassword:   "newHashedPassword",
			expectedError: false,
		},
		"invalid user ID": {
			userID:        999,
			newPassword:   "newHashedPassword",
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.UpdateUserPassword(tc.userID, tc.newPassword)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.users[tc.userID].Password != tc.newPassword {
					t.Fatalf("expected password to be updated, but it wasn't")
				}
			}
		})
	}
}
