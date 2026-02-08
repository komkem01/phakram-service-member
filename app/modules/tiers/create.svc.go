package tiers

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateTierService struct {
	NameTh       string          `json:"name_th"`
	NameEn       string          `json:"name_en"`
	MinSpending  decimal.Decimal `json:"min_spending"`
	IsActive     bool            `json:"is_active"`
	DiscountRate decimal.Decimal `json:"discount_rate"`
}

func (s *Service) CreateTierService(ctx context.Context, req *CreateTierService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`tiers.svc.create.start`)

	tier := &ent.TierEntity{
		ID:           uuid.New(),
		NameTh:       req.NameTh,
		NameEn:       req.NameEn,
		MinSpending:  req.MinSpending,
		IsActive:     req.IsActive,
		DiscountRate: req.DiscountRate,
	}
	if err := s.db.CreateTier(ctx, tier); err != nil {
		return err
	}
	span.AddEvent(`tiers.svc.create.success`)
	return nil
}
