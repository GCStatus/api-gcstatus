package ports

type UpdateSocialsRequest struct {
	Share     *bool   `json:"share" binding:"required"`
	Phone     *string `json:"phone,omitempty"`
	Github    *string `json:"github,omitempty"`
	Twitch    *string `json:"twitch,omitempty"`
	Twitter   *string `json:"twitter,omitempty"`
	Youtube   *string `json:"youtube,omitempty"`
	Facebook  *string `json:"facebook,omitempty"`
	Instagram *string `json:"instagram,omitempty"`
}

type ProfileRepository interface {
	UpdateSocials(profileID uint, request UpdateSocialsRequest) error
	UpdatePicture(profileID uint, path string) error
}
