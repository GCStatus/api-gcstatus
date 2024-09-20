package resources

import (
	"gcstatus/internal/domain"
)

// ProfileResource defines the structure of the profile response
type ProfileResource struct {
	ID        uint   `json:"id"`
	Share     bool   `json:"share"`
	Photo     string `json:"photo,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	Youtube   string `json:"youtube,omitempty"`
	Twitch    string `json:"twitch,omitempty"`
	Github    string `json:"github,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// TransformProfile is a function to transform a single profile
func TransformProfile(profile domain.Profile) *ProfileResource {
	return &ProfileResource{
		ID:        profile.ID,
		Share:     profile.Share,
		Photo:     profile.Photo,
		Phone:     profile.Phone,
		Facebook:  profile.Facebook,
		Instagram: profile.Instagram,
		Twitter:   profile.Twitter,
		Youtube:   profile.Youtube,
		Twitch:    profile.Twitch,
		Github:    profile.Github,
		CreatedAt: profile.CreatedAt.Format("2006-01-02T15:04:05"),
		UpdatedAt: profile.UpdatedAt.Format("2006-01-02T15:04:05"),
	}
}
