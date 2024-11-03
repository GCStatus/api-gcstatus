package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"
	"strings"
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

func (m *MockGameRepository) ExistsForStore(storeID uint, appID string) (bool, error) {
	for _, game := range m.games {
		for _, gameStore := range game.Stores {
			if gameStore.StoreID == storeID && gameStore.StoreGameID == appID {
				return true, nil
			}
		}
	}
	return false, nil
}

func (m *MockGameRepository) Search(input string) ([]domain.Game, error) {
	var games []domain.Game

	for _, game := range m.games {
		if strings.Contains(game.Title, input) || strings.Contains(game.ShortDescription, input) || strings.Contains(game.Description, input) || strings.Contains(game.Title, input) {
			games = append(games, *game)
		} else {
			return nil, errors.New("no one games found")
		}
	}

	return games, nil
}

func (m *MockGameRepository) FindByClassification(classification string, filterable string) ([]domain.Game, error) {
	var games []domain.Game

	for _, game := range m.games {
		switch classification {
		case "genres":
			for _, genre := range game.Genres {
				if genre.Genre.Slug == filterable {
					games = append(games, *game)
				}
			}
		case "categories":
			for _, category := range game.Categories {
				if category.Category.Slug == filterable {
					games = append(games, *game)
				}
			}
		case "tags":
			for _, tag := range game.Tags {
				if tag.Tag.Slug == filterable {
					games = append(games, *game)
				}
			}
		case "platforms":
			for _, platform := range game.Platforms {
				if platform.Platform.Slug == filterable {
					games = append(games, *game)
				}
			}
		default:
			return nil, errors.New("invalid classification")
		}
	}

	if len(games) == 0 {
		return nil, errors.New("games not found for this classification")
	}

	return games, nil
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

func TestMockGameRepository_Search(t *testing.T) {
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
		input       string
		expectError bool
	}{
		"valid game slug": {
			input:       "Game Test",
			expectError: false,
		},
		"invalid game slug": {
			input:       "invalid",
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			games, err := mockRepo.Search(tc.input)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if games != nil {
					t.Fatalf("expected nil games, got %v", games)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
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

func TestMockGameRepository_ExistsForStore(t *testing.T) {
	mockRepo := NewMockGameRepository()

	if err := mockRepo.CreateGame(&domain.Game{ID: 1, Slug: "game1", Stores: []domain.GameStore{
		{
			StoreID:     1,
			StoreGameID: "100",
		},
	}}); err != nil {
		t.Fatalf("Failed to create game")
	}
	if err := mockRepo.CreateGame(&domain.Game{ID: 2, Slug: "game2", Stores: []domain.GameStore{
		{
			StoreID:     2,
			StoreGameID: "200",
		},
	}}); err != nil {
		t.Fatalf("Failed to create game")
	}

	tests := map[string]struct {
		appID       string
		storeID     uint
		expected    bool
		expectError bool
	}{
		"game exists in store": {
			appID:       "100",
			storeID:     1,
			expected:    true,
			expectError: false,
		},
		"game does not exist in store": {
			appID:       "300",
			storeID:     1,
			expected:    false,
			expectError: false,
		},
		"appID does not match any game": {
			appID:       "100",
			storeID:     999,
			expected:    false,
			expectError: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			exists, err := mockRepo.ExistsForStore(tt.storeID, tt.appID)

			if tt.expectError && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if exists != tt.expected {
				t.Errorf("expected exists to be %v, got %v", tt.expected, exists)
			}
		})
	}
}

func TestMockGameRepository_FindByClassification(t *testing.T) {
	mockRepo := NewMockGameRepository()

	if err := mockRepo.CreateGame(&domain.Game{ID: 1, Genres: []domain.Genreable{
		{ID: 1, GenreableID: 1, GenreableType: "games", Genre: domain.Genre{ID: 1, Name: "Action", Slug: "action"}},
	}}); err != nil {
		t.Fatalf("Failed to create game: %+v", err)
	}
	if err := mockRepo.CreateGame(&domain.Game{ID: 2, Genres: []domain.Genreable{
		{ID: 2, GenreableID: 2, GenreableType: "games", Genre: domain.Genre{ID: 1, Name: "Action", Slug: "action"}},
	}, Platforms: []domain.Platformable{
		{ID: 1, PlatformableID: 1, PlatformableType: "games", Platform: domain.Platform{ID: 1, Name: "PS5", Slug: "ps5"}},
	}}); err != nil {
		t.Fatalf("Failed to create game: %+v", err)
	}
	if err := mockRepo.CreateGame(&domain.Game{ID: 3, Genres: []domain.Genreable{
		{ID: 3, GenreableID: 3, GenreableType: "games", Genre: domain.Genre{ID: 1, Name: "Action", Slug: "action"}},
	}, Platforms: []domain.Platformable{
		{ID: 2, PlatformableID: 3, PlatformableType: "games", Platform: domain.Platform{ID: 1, Name: "PS5", Slug: "ps5"}},
	}, Categories: []domain.Categoriable{
		{ID: 1, CategoriableID: 3, CategoriableType: "games", Category: domain.Category{ID: 1, Name: "Adventure", Slug: "adventure"}},
	}}); err != nil {
		t.Fatalf("Failed to create game: %+v", err)
	}

	tests := map[string]struct {
		classification string
		filter         string
		expectedCount  int
		expectError    bool
	}{
		"Find games by genre 'action'": {
			classification: "genres",
			filter:         "action",
			expectedCount:  3,
			expectError:    false,
		},
		"Find games by platform 'ps5'": {
			classification: "platforms",
			filter:         "ps5",
			expectedCount:  2,
			expectError:    false,
		},
		"Find games by category 'adventure'": {
			classification: "categories",
			filter:         "adventure",
			expectedCount:  1,
			expectError:    false,
		},
		"No games for nonexistent tag": {
			classification: "tags",
			filter:         "nonexistent",
			expectedCount:  0,
			expectError:    true,
		},
		"Invalid classification": {
			classification: "invalid_classification",
			filter:         "anything",
			expectedCount:  0,
			expectError:    true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			games, err := mockRepo.FindByClassification(tt.classification, tt.filter)

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
