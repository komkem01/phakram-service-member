package tiers

import (
	"context"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UpdateTierService struct {
	NameTh       string           `json:"name_th"`
	NameEn       string           `json:"name_en"`
	MinSpending  *decimal.Decimal `json:"min_spending"`
	IsActive     *bool            `json:"is_active"`
	DiscountRate *decimal.Decimal `json:"discount_rate"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateTierService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`tiers.svc.update.start`)

	data, err := s.db.GetTierByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.NameTh != "" {
		data.NameTh = req.NameTh
	}
	if req.NameEn != "" {
		data.NameEn = req.NameEn
	}
	if req.MinSpending != nil {
		data.MinSpending = *req.MinSpending
	}
	if req.IsActive != nil {
		data.IsActive = *req.IsActive
	}
	if req.DiscountRate != nil {
		data.DiscountRate = *req.DiscountRate
	}

	if err := s.db.UpdateTier(ctx, data); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`tiers.svc.update.success`)
	return nil
}
