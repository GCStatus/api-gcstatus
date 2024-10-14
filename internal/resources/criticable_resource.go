package resources

import (
	"gcstatus/internal/domain"
	"gcstatus/internal/utils"

	"github.com/shopspring/decimal"
)

type CriticableResource struct {
	ID       uint            `json:"id"`
	URL      string          `json:"url"`
	Rate     decimal.Decimal `json:"rate"`
	PostedAt string          `json:"posted_at"`
	Critic   CriticResource  `json:"critic"`
}

func TransformCriticable(criticable domain.Criticable) CriticableResource {
	var postedAt string
	if !criticable.PostedAt.IsZero() {
		postedAt = utils.FormatTimestamp(criticable.PostedAt)
	}

	return CriticableResource{
		ID:       criticable.ID,
		URL:      criticable.URL,
		Rate:     criticable.Rate,
		PostedAt: postedAt,
		Critic:   TransformCritic(criticable.Critic),
	}
}
