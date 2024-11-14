package test_mocks

import (
	"gcstatus/internal/domain"
	"testing"

	"gorm.io/gorm"
)

func CreateDummyComment(t *testing.T, dbConn *gorm.DB, overrides *domain.Commentable) (*domain.Commentable, error) {
	defaultComment := domain.Commentable{
		Comment:         "Testing comment",
		CommentableID:   1,
		CommentableType: "games",
	}

	if overrides != nil {
		if overrides.Comment != "" {
			defaultComment.Comment = overrides.Comment
		}
		if overrides.CommentableID != 0 {
			defaultComment.CommentableID = overrides.CommentableID
		}
		if overrides.CommentableType != "" {
			defaultComment.CommentableType = overrides.CommentableType
		}
		if overrides.User.ID != 0 {
			defaultComment.User = overrides.User
		} else {
			user, err := CreateDummyUser(t, dbConn, &overrides.User)
			if err != nil {
				t.Fatalf("failed to create dummy user for comment: %+v", err)
			}

			defaultComment.User = *user
		}
	}

	if err := dbConn.Create(&defaultComment).Error; err != nil {
		return nil, err
	}

	return &defaultComment, nil
}
