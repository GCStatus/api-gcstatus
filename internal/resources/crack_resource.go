package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/pkg/utils"
)

type CrackResource struct {
	ID         uint                `json:"id"`
	Status     string              `json:"status"`
	CrackedAt  *string             `json:"cracked_at"`
	By         *CrackerResource    `json:"by"`
	Protection *ProtectionResource `json:"protection"`
}

func TransformCrack(crack *domain.Crack) *CrackResource {
	resource := CrackResource{
		ID:     crack.ID,
		Status: crack.Status,
	}

	if crack.CrackedAt != nil {
		formattedTime := utils.FormatTimestamp(*crack.CrackedAt)
		resource.CrackedAt = &formattedTime
	}

	if crack.Cracker.ID != 0 {
		resource.By = TransformCracker(crack.Cracker)
	} else {
		resource.By = nil
	}

	if crack.Protection.ID != 0 {
		resource.Protection = TransformProtection(crack.Protection)
	} else {
		resource.Protection = nil
	}

	return &resource
}