package tests

import (
	"errors"
	"gcstatus/internal/domain"
	"testing"
	"time"
)

type MockCommentRepository struct {
	comments map[uint]*domain.Commentable
}

func NewMockCommentRepository() *MockCommentRepository {
	return &MockCommentRepository{
		comments: make(map[uint]*domain.Commentable),
	}
}

func (m *MockCommentRepository) Create(commentable *domain.Commentable) error {
	if commentable == nil {
		return errors.New("invalid comment data")
	}
	m.comments[commentable.ID] = commentable
	return nil
}

func (m *MockCommentRepository) FindByID(id uint) (*domain.Commentable, error) {
	for _, comment := range m.comments {
		if comment.ID == id {
			return comment, nil
		}
	}

	return nil, errors.New("comment not found")
}

func (m *MockCommentRepository) Delete(id uint) error {
	if _, exists := m.comments[id]; !exists {
		return errors.New("commentable not found")
	}
	delete(m.comments, id)
	return nil
}

func TestMockCommentRepository_Create(t *testing.T) {
	mockRepo := NewMockCommentRepository()

	testCases := map[string]struct {
		input         *domain.Commentable
		expectedError bool
	}{
		"valid input": {
			input: &domain.Commentable{
				CommentableID:   1,
				CommentableType: "games",
				UserID:          1,
			},
			expectedError: false,
		},
		"nil input": {
			input:         nil,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.Create(tc.input)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.comments[tc.input.ID] == nil {
					t.Fatalf("expected commentable to be created, but it wasn't")
				}
			}
		})
	}
}

func TestMockCommentRepository_FindById(t *testing.T) {
	fixedTime := time.Now()

	mockRepo := NewMockCommentRepository()
	if err := mockRepo.Create(&domain.Commentable{
		ID:              1,
		CommentableID:   1,
		CommentableType: "games",
		CreatedAt:       fixedTime,
		UpdatedAt:       fixedTime,
	}); err != nil {
		t.Fatalf("failed to create the comment: %s", err.Error())
	}

	testCases := map[string]struct {
		commentID   uint
		expectError bool
	}{
		"valid comment ID": {
			commentID:   1,
			expectError: false,
		},
		"invalid comment ID": {
			commentID:   999,
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			comment, err := mockRepo.FindByID(tc.commentID)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if comment != nil {
					t.Fatalf("expected nil comment, got %v", comment)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if comment == nil || comment.ID != tc.commentID {
					t.Fatalf("expected comment ID %d, got %v", tc.commentID, comment)
				}
			}
		})
	}
}

func TestMockCommentRepository_Delete(t *testing.T) {
	mockRepo := NewMockCommentRepository()

	if err := mockRepo.Create(&domain.Commentable{
		ID:              1,
		CommentableID:   1,
		CommentableType: "games",
		UserID:          1,
	}); err != nil {
		t.Fatalf("failed to create the commentable: %s", err.Error())
	}

	testCases := map[string]struct {
		id            uint
		expectedError bool
	}{
		"valid ID": {
			id:            1,
			expectedError: false,
		},
		"invalid ID": {
			id:            999,
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := mockRepo.Delete(tc.id)

			if tc.expectedError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if mockRepo.comments[tc.id] != nil {
					t.Fatalf("expected heartable to be deleted, but it wasn't")
				}
			}
		})
	}
}
