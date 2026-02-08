package tiers

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InfoTierServiceResponses struct {
	ID           uuid.UUID `json:"id"`
	NameTh       string          `json:"name_th"`
	NameEn       string          `json:"name_en"`
	MinSpending  decimal.Decimal `json:"min_spending"`
	IsActive     bool            `json:"is_active"`
	DiscountRate decimal.Decimal `json:"discount_rate"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoTierServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`tiers.svc.info.start`)

	data, err := s.db.GetTierByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoTierServiceResponses{
		ID:           data.ID,
		NameTh:       data.NameTh,
		NameEn:       data.NameEn,
		MinSpending:  data.MinSpending,
		IsActive:     data.IsActive,
		DiscountRate: data.DiscountRate,
		CreatedAt:    data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`tiers.svc.info.success`)
	return resp, nil
}
