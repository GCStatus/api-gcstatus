package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"testing"
	"time"
)

type MockGameRepository struct {
	games map[uint]*domain.Game
}

func NewMockGameRepository() *MockGameRepository {
	return &MockGameRepository{
		games: make(map[uint]*domain.Game),
	}
}

func (m *MockGameRepository) GetAll() ([]domain.Game, error) {
	var games []domain.Game
	for _, game := range m.games {
		games = append(games, *game)
	}
	return games, nil
}

func (m *MockGameRepository) FindByID(id uint) (*domain.Game, error) {
	for _, game := range m.games {
		if game.ID == id {
			return game, nil
		}
	}
	return nil, errors.New("game not found")
}

func (m *MockGameRepository) CreateGame(game *domain.Game) error {
	if game == nil {
		return errors.New("invalid game data")
	}
	m.games[game.ID] = game
	return nil
}

func TestMockGameRepository_GetAll(t *testing.T) {
	fixedTime := time.Now()

	testCases := map[string]struct {
		expectedGameCount int
		mockCreateGames   func(repo *MockGameRepository)
	}{
		"multiple levels": {
			expectedGameCount: 2,
			mockCreateGames: func(repo *MockGameRepository) {
				if err := repo.CreateGame(&domain.Game{
					ID:               1,
					Slug:             "valid-1",
					Age:              18,
					Title:            "Game Test",
					Cover:            "https://placehold.co/600x400/EEE/31343C",
					About:            "About game",
					Description:      "Description",
					ShortDescription: "Short description",
					Free:             false,
					ReleaseDate:      fixedTime,
					CreatedAt:        fixedTime,
					UpdatedAt:        fixedTime,
				}); err != nil {
					t.Fatalf("failed to create the game: %s", err.Error())
				}
				if err := repo.CreateGame(&domain.Game{
					ID:               2,
					Slug:             "valid-2",
					Age:              18,
					Title:            "Game Test",
					Cover:            "https://placehold.co/600x400/EEE/31343C",
					About:            "About game",
					Description:      "Description",
					ShortDescription: "Short description",
					Free:             false,
					ReleaseDate:      fixedTime,
					CreatedAt:        fixedTime,
					UpdatedAt:        fixedTime,
				}); err != nil {
					t.Fatalf("failed to create the game: %s", err.Error())
				}
			},
		},
		"no levels": {
			expectedGameCount: 0,
			mockCreateGames:   func(repo *MockGameRepository) {},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockRepo := NewMockGameRepository()

			tc.mockCreateGames(mockRepo)

			levels, err := mockRepo.GetAll()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(levels) != tc.expectedGameCount {
				t.Fatalf("expected %d levels, got %d", tc.expectedGameCount, len(levels))
			}
		})
	}
}

func TestMockGameRepository_FindByID(t *testing.T) {
	fixedTime := time.Now()

	mockRepo := NewMockGameRepository()
	if err := mockRepo.CreateGame(&domain.Game{
		ID:               1,
		Slug:             "valid",
		Age:              18,
		Title:            "Game Test",
		Cover:            "https://placehold.co/600x400/EEE/31343C",
		About:            "About game",
		Description:      "Description",
		ShortDescription: "Short description",
		Free:             false,
		ReleaseDate:      fixedTime,
		CreatedAt:        fixedTime,
		UpdatedAt:        fixedTime,
	}); err != nil {
		t.Fatalf("failed to create the game: %s", err.Error())
	}

	testCases := map[string]struct {
		gameID      uint
		expectError bool
	}{
		"valid game id": {
			gameID:      1,
			expectError: false,
		},
		"invalid game id": {
			gameID:      2,
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			game, err := mockRepo.FindByID(tc.gameID)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if game != nil {
					t.Fatalf("expected nil game, got %v", game)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if game == nil || game.ID != tc.gameID {
					t.Fatalf("expected game Slug %d, got %v", tc.gameID, game)
				}
			}
		})
	}
}
