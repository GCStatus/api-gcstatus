package tests

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/resources"
	"gcstatus/pkg/utils"
	"testing"
	"time"
)

func TestTransformCommentable(t *testing.T) {
	fixedTime := time.Now()
	formattedTime := utils.FormatTimestamp(fixedTime)

	testCases := map[string]struct {
		userID   uint
		input    domain.Commentable
		expected resources.CommentableResource
	}{
		"liked reply": {
			input: domain.Commentable{
				ID:        1,
				Comment:   "Fake comment",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				User: domain.User{
					ID:        1,
					Name:      "John Doe",
					Email:     "johndoe@example.com",
					Nickname:  "johnny",
					CreatedAt: fixedTime,
					Profile: domain.Profile{
						Share: true,
						Photo: "photo-key-1",
					},
				},
				Replies: []domain.Commentable{
					{
						ID:        1,
						Comment:   "Fake comment",
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						User: domain.User{
							ID:        1,
							Name:      "John Doe",
							Email:     "johndoe@example.com",
							Nickname:  "johnny",
							CreatedAt: fixedTime,
							Profile: domain.Profile{
								Share: true,
								Photo: "photo-key-1",
							},
						},
						Hearts: []domain.Heartable{
							{
								ID:            1,
								HeartableID:   1,
								HeartableType: "commentables",
								UserID:        1,
							},
						},
					},
				},
			},
			expected: resources.CommentableResource{
				ID:          1,
				Comment:     "Fake comment",
				CreatedAt:   formattedTime,
				HeartsCount: 0,
				IsHearted:   false,
				UpdatedAt:   formattedTime,
				By: resources.MinimalUserResource{
					ID:        1,
					Name:      "John Doe",
					Photo:     utils.StringPtr("https://mock-presigned-url.com/photo-key-1"),
					Email:     "johndoe@example.com",
					Nickname:  "johnny",
					CreatedAt: formattedTime,
				},
				Replies: []resources.CommentableResource{
					{
						ID:          1,
						Comment:     "Fake comment",
						CreatedAt:   formattedTime,
						HeartsCount: 1,
						IsHearted:   true,
						UpdatedAt:   formattedTime,
						By: resources.MinimalUserResource{
							ID:        1,
							Name:      "John Doe",
							Photo:     utils.StringPtr("https://mock-presigned-url.com/photo-key-1"),
							Email:     "johndoe@example.com",
							Nickname:  "johnny",
							CreatedAt: formattedTime,
						},
					},
				},
			},
		},
		"not liked comment": {
			input: domain.Commentable{
				ID:        1,
				Comment:   "Fake comment",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				User: domain.User{
					ID:        1,
					Name:      "John Doe",
					Email:     "johndoe@example.com",
					Nickname:  "johnny",
					CreatedAt: fixedTime,
					Profile: domain.Profile{
						Share: true,
						Photo: "photo-key-1",
					},
				},
				Hearts: []domain.Heartable{
					{
						ID:            1,
						HeartableID:   1,
						HeartableType: "commentables",
						UserID:        2,
					},
				},
			},
			expected: resources.CommentableResource{
				ID:          1,
				Comment:     "Fake comment",
				CreatedAt:   formattedTime,
				HeartsCount: 1,
				IsHearted:   false,
				UpdatedAt:   formattedTime,
				By: resources.MinimalUserResource{
					ID:        1,
					Name:      "John Doe",
					Photo:     utils.StringPtr("https://mock-presigned-url.com/photo-key-1"),
					Email:     "johndoe@example.com",
					Nickname:  "johnny",
					CreatedAt: formattedTime,
				},
				Replies: []resources.CommentableResource{},
			},
		},
		"liked comment": {
			input: domain.Commentable{
				ID:        1,
				Comment:   "Fake comment",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				User: domain.User{
					ID:        1,
					Name:      "John Doe",
					Email:     "johndoe@example.com",
					Nickname:  "johnny",
					CreatedAt: fixedTime,
					Profile: domain.Profile{
						Share: true,
						Photo: "photo-key-1",
					},
				},
				Hearts: []domain.Heartable{
					{
						ID:            1,
						HeartableID:   1,
						HeartableType: "commentables",
						UserID:        1,
					},
				},
			},
			expected: resources.CommentableResource{
				ID:          1,
				Comment:     "Fake comment",
				CreatedAt:   formattedTime,
				HeartsCount: 1,
				IsHearted:   true,
				UpdatedAt:   formattedTime,
				By: resources.MinimalUserResource{
					ID:        1,
					Name:      "John Doe",
					Photo:     utils.StringPtr("https://mock-presigned-url.com/photo-key-1"),
					Email:     "johndoe@example.com",
					Nickname:  "johnny",
					CreatedAt: formattedTime,
				},
				Replies: []resources.CommentableResource{},
			},
		},
		"without replies": {
			input: domain.Commentable{
				ID:        1,
				Comment:   "Fake comment",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				User: domain.User{
					ID:        1,
					Name:      "John Doe",
					Email:     "johndoe@example.com",
					Nickname:  "johnny",
					CreatedAt: fixedTime,
					Profile: domain.Profile{
						Share: true,
						Photo: "photo-key-1",
					},
				},
			},
			expected: resources.CommentableResource{
				ID:        1,
				Comment:   "Fake comment",
				CreatedAt: formattedTime,
				UpdatedAt: formattedTime,
				By: resources.MinimalUserResource{
					ID:        1,
					Name:      "John Doe",
					Photo:     utils.StringPtr("https://mock-presigned-url.com/photo-key-1"),
					Email:     "johndoe@example.com",
					Nickname:  "johnny",
					CreatedAt: formattedTime,
				},
				Replies: []resources.CommentableResource{},
			},
		},
		"with replies": {
			input: domain.Commentable{
				ID:        1,
				Comment:   "Main comment",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				User: domain.User{
					ID:        1,
					Name:      "Main User",
					Email:     "mainuser@example.com",
					Nickname:  "mainuser",
					CreatedAt: fixedTime,
				},
				Replies: []domain.Commentable{
					{
						ID:        2,
						Comment:   "Reply to main comment",
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						User: domain.User{
							ID:        2,
							Name:      "Reply User",
							Email:     "replyuser@example.com",
							Nickname:  "replyuser",
							CreatedAt: fixedTime,
							Profile: domain.Profile{
								Share: true,
							},
						},
					},
				},
			},
			expected: resources.CommentableResource{
				ID:        1,
				Comment:   "Main comment",
				CreatedAt: formattedTime,
				UpdatedAt: formattedTime,
				By: resources.MinimalUserResource{
					ID:        1,
					Name:      "Main User",
					Photo:     nil,
					Email:     "mainuser@example.com",
					Nickname:  "mainuser",
					CreatedAt: formattedTime,
				},
				Replies: []resources.CommentableResource{
					{
						ID:        2,
						Comment:   "Reply to main comment",
						CreatedAt: formattedTime,
						UpdatedAt: formattedTime,
						By: resources.MinimalUserResource{
							ID:        2,
							Name:      "Reply User",
							Photo:     nil,
							Email:     "replyuser@example.com",
							Nickname:  "replyuser",
							CreatedAt: formattedTime,
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockS3Client := &MockS3Client{}
			result := resources.TransformCommentable(tc.input, mockS3Client, tc.userID)

			if !compareCommentableResources(tc.expected, result) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func compareCommentableResources(expected, actual resources.CommentableResource) bool {
	if expected.ID != actual.ID ||
		expected.Comment != actual.Comment ||
		expected.CreatedAt != actual.CreatedAt ||
		expected.UpdatedAt != actual.UpdatedAt {
		return false
	}

	if !compareMinimalUserResources(expected.By, actual.By) {
		return false
	}

	if len(expected.Replies) != len(actual.Replies) {
		return false
	}

	for i := range expected.Replies {
		if !compareCommentableResources(expected.Replies[i], actual.Replies[i]) {
			return false
		}
	}

	return true
}

func compareMinimalUserResources(expected, actual resources.MinimalUserResource) bool {
	if expected.ID != actual.ID ||
		expected.Name != actual.Name ||
		expected.Email != actual.Email ||
		expected.Nickname != actual.Nickname ||
		expected.CreatedAt != actual.CreatedAt {
		return false
	}

	if (expected.Photo == nil && actual.Photo != nil) ||
		(expected.Photo != nil && actual.Photo == nil) {
		return false
	}

	if expected.Photo != nil && actual.Photo != nil && *expected.Photo != *actual.Photo {
		return false
	}

	return true
}
