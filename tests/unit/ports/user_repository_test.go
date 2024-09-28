package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"gcstatus/internal/ports"
	"gcstatus/pkg/utils"
	"testing"
	"time"
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

func (m *MockUserRepository) UpdateUserNickAndEmail(userID uint, request ports.UpdateNickAndEmailRequest) error {
	user, exists := m.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	if user.Password != request.Password {
		return errors.New("password does not match")
	}

	for id, u := range m.users {
		if id != userID {
			if u.Email == request.Email {
				return errors.New("email already in use")
			}
			if u.Nickname == request.Nickname {
				return errors.New("nickname already in use")
			}
		}
	}

	user.Nickname = request.Nickname
	user.Email = request.Email

	m.users[userID] = user
	return nil
}

func (m *MockUserRepository) UpdateUserBasics(userID uint, request ports.UpdateUserBasicsRequest) error {
	user, exists := m.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	birthdate, err := time.Parse("2006-01-02T15:04:05", request.Birthdate)
	if err != nil {
		birthdate, err = time.Parse("2006-01-02", request.Birthdate)
		if err != nil {
			return errors.New("birthdate incorrectly formatted")
		}
	}

	if time.Since(birthdate).Hours() < 14*365*24 {
		return errors.New("user lower than 14 years")
	}

	user.Name = request.Name
	user.Birthdate = birthdate
	m.users[userID] = user

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

func TestMockUserRepository_UpdateUserNickAndEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()

	err := mockRepo.CreateUser(&domain.User{
		ID:       1,
		Email:    "user@example.com",
		Nickname: "user1",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("failed to create the user: %s", err.Error())
	}

	err = mockRepo.CreateUser(&domain.User{
		ID:       2,
		Email:    "duplicate@example.com",
		Nickname: "duplicateNick",
		Password: "password456",
	})
	if err != nil {
		t.Fatalf("failed to create the second user: %s", err.Error())
	}

	testCases := map[string]struct {
		userID        uint
		newEmail      string
		newNickname   string
		password      string
		expectedError bool
		expectedMsg   string
	}{
		"success": {
			userID:        1,
			newNickname:   "user2",
			newEmail:      "user2@example.com",
			password:      "password123",
			expectedError: false,
		},
		"invalid user ID": {
			userID:        999,
			newNickname:   "user2",
			newEmail:      "user2@example.com",
			password:      "password123",
			expectedError: true,
			expectedMsg:   "user not found",
		},
		"duplicated email": {
			userID:        1,
			newNickname:   "user2",
			newEmail:      "duplicate@example.com",
			password:      "password123",
			expectedError: true,
			expectedMsg:   "email already in use",
		},
		"duplicated nickname": {
			userID:        1,
			newNickname:   "duplicateNick",
			newEmail:      "user2@example.com",
			password:      "password123",
			expectedError: true,
			expectedMsg:   "nickname already in use",
		},
		"password does not match": {
			userID:        1,
			newNickname:   "user2",
			newEmail:      "user2@example.com",
			password:      "wrongpassword",
			expectedError: true,
			expectedMsg:   "password does not match",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := ports.UpdateNickAndEmailRequest{
				Password: tc.password,
				Nickname: tc.newNickname,
				Email:    tc.newEmail,
			}

			err := mockRepo.UpdateUserNickAndEmail(tc.userID, request)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if err.Error() != tc.expectedMsg {
					t.Fatalf("expected error message: %s, got: %s", tc.expectedMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				updatedUser := mockRepo.users[tc.userID]
				if updatedUser.Email != tc.newEmail {
					t.Fatalf("expected email to be updated to %s, but got %s", tc.newEmail, updatedUser.Email)
				}
				if updatedUser.Nickname != tc.newNickname {
					t.Fatalf("expected nickname to be updated to %s, but got %s", tc.newNickname, updatedUser.Nickname)
				}
			}
		})
	}
}

func TestMockUserRepository_UpdateUserBasics(t *testing.T) {
	mockRepo := NewMockUserRepository()
	fixedTime := time.Now()

	err := mockRepo.CreateUser(&domain.User{
		ID:        1,
		Name:      "User",
		Birthdate: fixedTime,
	})
	if err != nil {
		t.Fatalf("failed to create the user: %s", err.Error())
	}

	testCases := map[string]struct {
		userID        uint
		newName       string
		newBirthdate  string
		expectedError bool
		expectedMsg   string
	}{
		"success": {
			userID:        1,
			newName:       "User 2",
			newBirthdate:  "2000-01-01T00:00:00",
			expectedError: false,
		},
		"invalid user ID": {
			userID:        999,
			newName:       "User 2",
			newBirthdate:  "2000-01-01T00:00:00",
			expectedError: true,
			expectedMsg:   "user not found",
		},
		"less than 14": {
			userID:        1,
			newName:       "User 2",
			newBirthdate:  utils.FormatTimestamp(time.Now()),
			expectedError: true,
			expectedMsg:   "user lower than 14 years",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := ports.UpdateUserBasicsRequest{
				Name:      tc.newName,
				Birthdate: tc.newBirthdate,
			}

			err := mockRepo.UpdateUserBasics(tc.userID, request)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if err.Error() != tc.expectedMsg {
					t.Fatalf("expected error message: %s, got: %s", tc.expectedMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				updatedUser := mockRepo.users[tc.userID]
				if updatedUser.Name != tc.newName {
					t.Fatalf("expected name to be updated to %s, but got %s", tc.newName, updatedUser.Name)
				}
				if utils.FormatTimestamp(updatedUser.Birthdate) != tc.newBirthdate {
					t.Fatalf("expected birthdate to be updated to %s, but got %s", tc.newBirthdate, utils.FormatTimestamp(updatedUser.Birthdate))
				}
			}
		})
	}
}
