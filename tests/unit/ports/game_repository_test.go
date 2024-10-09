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
