package usecases

import (
	"gcstatus/internal/ports"
)

type ProfileService struct {
	repo ports.ProfileRepository
}

func NewProfileService(repo ports.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (h *ProfileService) UpdateSocials(profileID uint, request ports.UpdateSocialsRequest) error {
	return h.repo.UpdateSocials(profileID, request)
}

func (h *ProfileService) UpdatePicture(profileID uint, path string) error {
	return h.repo.UpdatePicture(profileID, path)
}
