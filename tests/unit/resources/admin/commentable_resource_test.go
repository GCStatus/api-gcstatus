package tests

import (
	"gcstatus/internal/domain"
	resources_admin "gcstatus/internal/resources/admin"
	"gcstatus/internal/utils"
	testutils "gcstatus/tests/utils"
	"testing"
	"time"
)

func TestTransformCommentable(t *testing.T) {
	fixedTime := time.Now()
	formattedTime := utils.FormatTimestamp(fixedTime)

	testCases := map[string]struct {
		input    domain.Commentable
		expected resources_admin.CommentableResource
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
			expected: resources_admin.CommentableResource{
				ID:          1,
				Comment:     "Fake comment",
				CreatedAt:   formattedTime,
				HeartsCount: 0,
				UpdatedAt:   formattedTime,
				By: resources_admin.MinimalUserResource{
					ID:        1,
					Name:      "John Doe",
					Photo:     utils.StringPtr("https://mock-presigned-url.com/photo-key-1"),
					Email:     "johndoe@example.com",
					Nickname:  "johnny",
					CreatedAt: formattedTime,
				},
				Replies: []resources_admin.CommentableResource{
					{
						ID:          1,
						Comment:     "Fake comment",
						CreatedAt:   formattedTime,
						HeartsCount: 1,
						UpdatedAt:   formattedTime,
						By: resources_admin.MinimalUserResource{
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
			expected: resources_admin.CommentableResource{
				ID:          1,
				Comment:     "Fake comment",
				CreatedAt:   formattedTime,
				HeartsCount: 1,
				UpdatedAt:   formattedTime,
				By: resources_admin.MinimalUserResource{
					ID:        1,
					Name:      "John Doe",
					Photo:     utils.StringPtr("https://mock-presigned-url.com/photo-key-1"),
					Email:     "johndoe@example.com",
					Nickname:  "johnny",
					CreatedAt: formattedTime,
				},
				Replies: []resources_admin.CommentableResource{},
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
			expected: resources_admin.CommentableResource{
				ID:          1,
				Comment:     "Fake comment",
				CreatedAt:   formattedTime,
				HeartsCount: 1,
				UpdatedAt:   formattedTime,
				By: resources_admin.MinimalUserResource{
					ID:        1,
					Name:      "John Doe",
					Photo:     utils.StringPtr("https://mock-presigned-url.com/photo-key-1"),
					Email:     "johndoe@example.com",
					Nickname:  "johnny",
					CreatedAt: formattedTime,
				},
				Replies: []resources_admin.CommentableResource{},
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
			expected: resources_admin.CommentableResource{
				ID:        1,
				Comment:   "Fake comment",
				CreatedAt: formattedTime,
				UpdatedAt: formattedTime,
				By: resources_admin.MinimalUserResource{
					ID:        1,
					Name:      "John Doe",
					Photo:     utils.StringPtr("https://mock-presigned-url.com/photo-key-1"),
					Email:     "johndoe@example.com",
					Nickname:  "johnny",
					CreatedAt: formattedTime,
				},
				Replies: []resources_admin.CommentableResource{},
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
			expected: resources_admin.CommentableResource{
				ID:        1,
				Comment:   "Main comment",
				CreatedAt: formattedTime,
				UpdatedAt: formattedTime,
				By: resources_admin.MinimalUserResource{
					ID:        1,
					Name:      "Main User",
					Photo:     nil,
					Email:     "mainuser@example.com",
					Nickname:  "mainuser",
					CreatedAt: formattedTime,
				},
				Replies: []resources_admin.CommentableResource{
					{
						ID:        2,
						Comment:   "Reply to main comment",
						CreatedAt: formattedTime,
						UpdatedAt: formattedTime,
						By: resources_admin.MinimalUserResource{
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
			mockS3Client := &testutils.MockS3Client{}
			result := resources_admin.TransformCommentable(tc.input, mockS3Client)

			if !compareCommentableResources_admin(tc.expected, result) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func compareCommentableResources_admin(expected, actual resources_admin.CommentableResource) bool {
	if expected.ID != actual.ID ||
		expected.Comment != actual.Comment ||
		expected.CreatedAt != actual.CreatedAt ||
		expected.UpdatedAt != actual.UpdatedAt {
		return false
	}

	if !compareMinimalUserResources_admin(expected.By, actual.By) {
		return false
	}

	if len(expected.Replies) != len(actual.Replies) {
		return false
	}

	for i := range expected.Replies {
		if !compareCommentableResources_admin(expected.Replies[i], actual.Replies[i]) {
			return false
		}
	}

	return true
}

func compareMinimalUserResources_admin(expected, actual resources_admin.MinimalUserResource) bool {
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
