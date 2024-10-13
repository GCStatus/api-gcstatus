package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
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
func (m *MockGameRepository) FindBySlug(slug string) (*domain.Game, error) {
	for _, game := range m.games {
		if game.Slug == slug {
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

func (m *MockGameRepository) HomeGames() ([]domain.Game, []domain.Game, []domain.Game, *domain.Game, []domain.Game, error) {
	var hotGames, popularGames, mostHeartedGames, upcomingGames []domain.Game
	var nextGreatReleaseGame *domain.Game

	for _, game := range m.games {
		switch game.Condition {
		case "hot":
			if len(hotGames) < 9 {
				hotGames = append(hotGames, *game)
			}
		case "popular":
			if len(popularGames) < 9 {
				popularGames = append(popularGames, *game)
			}
		}

		if len(mostHeartedGames) < 9 {
			mostHeartedGames = append(mostHeartedGames, *game)
		}

		if game.GreatRelease && game.ReleaseDate.After(time.Now()) && nextGreatReleaseGame == nil {
			nextGreatReleaseGame = game
		}

		if game.ReleaseDate.After(time.Now()) {
			upcomingGames = append(upcomingGames, *game)
		}
	}

	return hotGames, popularGames, mostHeartedGames, nextGreatReleaseGame, upcomingGames, nil
}

func (m *MockGameRepository) FindGamesByCondition(condition string, limit *uint) ([]domain.Game, error) {
	var games []domain.Game

	for _, game := range m.games {
		if game.Condition == condition {
			games = append(games, *game)
			if limit != nil && len(games) >= int(*limit) {
				break
			}
		}
	}

	return games, nil
}

func TestMockGameRepository_FindBySlug(t *testing.T) {
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
		t.Fatalf("failed to create the slug: %s", err.Error())
	}

	testCases := map[string]struct {
		gameSlug    string
		expectError bool
	}{
		"valid game slug": {
			gameSlug:    "valid",
			expectError: false,
		},
		"invalid game slug": {
			gameSlug:    "invalid",
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			game, err := mockRepo.FindBySlug(tc.gameSlug)

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
				if game == nil || game.Slug != tc.gameSlug {
					t.Fatalf("expected game Slug %s, got %v", tc.gameSlug, game)
				}
			}
		})
	}
}

func TestMockGameRepository_HomeGames(t *testing.T) {
	mockRepo := NewMockGameRepository()

	if err := mockRepo.CreateGame(&domain.Game{ID: 1, Slug: "game1", Condition: "hot", ReleaseDate: time.Now().AddDate(0, 0, 1), GreatRelease: true}); err != nil {
		t.Fatalf("Failed to create game")
	}
	if err := mockRepo.CreateGame(&domain.Game{ID: 2, Slug: "game2", Condition: "hot", ReleaseDate: time.Now().AddDate(0, 0, 1)}); err != nil {
		t.Fatalf("Failed to create game")
	}
	if err := mockRepo.CreateGame(&domain.Game{ID: 3, Slug: "game3", Condition: "popular", ReleaseDate: time.Now().AddDate(0, 0, 1)}); err != nil {
		t.Fatalf("Failed to create game")
	}
	if err := mockRepo.CreateGame(&domain.Game{ID: 4, Slug: "game4", Condition: "popular", ReleaseDate: time.Now().AddDate(0, 0, 1)}); err != nil {
		t.Fatalf("Failed to create game")
	}

	tests := []struct {
		name                     string
		expectedHotCount         int
		expectedPopularCount     int
		expectedMostHeartedCount int
		expectedUpcomingCount    int
		expectNextGreatRelease   bool
	}{
		{
			name:                     "default case",
			expectedHotCount:         2,
			expectedPopularCount:     2,
			expectedMostHeartedCount: 4,
			expectedUpcomingCount:    4,
			expectNextGreatRelease:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hotGames, popularGames, mostHeartedGames, nextGreatReleaseGame, upcomingGames, err := mockRepo.HomeGames()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(hotGames) != tt.expectedHotCount {
				t.Errorf("expected %d hot games, got %d", tt.expectedHotCount, len(hotGames))
			}

			if len(popularGames) != tt.expectedPopularCount {
				t.Errorf("expected %d popular games, got %d", tt.expectedPopularCount, len(popularGames))
			}

			if len(mostHeartedGames) != tt.expectedMostHeartedCount {
				t.Errorf("expected %d most hearted games, got %d", tt.expectedMostHeartedCount, len(mostHeartedGames))
			}

			if len(upcomingGames) != tt.expectedUpcomingCount {
				t.Errorf("expected %d upcoming games, got %d", tt.expectedUpcomingCount, len(upcomingGames))
			}

			if (nextGreatReleaseGame != nil) != tt.expectNextGreatRelease {
				t.Errorf("expected next great release presence %v, got %v", tt.expectNextGreatRelease, nextGreatReleaseGame != nil)
			}
		})
	}
}

func TestMockGameRepository_FindGamesByCondition(t *testing.T) {
	mockRepo := NewMockGameRepository()

	if err := mockRepo.CreateGame(&domain.Game{ID: 1, Slug: "game1", Condition: "hot"}); err != nil {
		t.Fatalf("Failed to create game")
	}
	if err := mockRepo.CreateGame(&domain.Game{ID: 2, Slug: "game2", Condition: "hot"}); err != nil {
		t.Fatalf("Failed to create game")
	}
	if err := mockRepo.CreateGame(&domain.Game{ID: 3, Slug: "game3", Condition: "popular"}); err != nil {
		t.Fatalf("Failed to create game")
	}

	tests := []struct {
		name          string
		condition     string
		limit         *uint
		expectedCount int
		expectError   bool
	}{
		{
			name:          "hot games with limit",
			condition:     "hot",
			limit:         utils.UintPtr(2),
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "popular games without limit",
			condition:     "popular",
			limit:         nil,
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "no games for nonexistent condition",
			condition:     "nonexistent",
			limit:         nil,
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			games, err := mockRepo.FindGamesByCondition(tt.condition, tt.limit)

			if tt.expectError && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(games) != tt.expectedCount {
				t.Errorf("expected %d games, got %d", tt.expectedCount, len(games))
			}
		})
	}
}
