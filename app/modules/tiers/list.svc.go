package tiers

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListTierServiceRequest struct {
	base.RequestPaginate
}

type ListTierServiceResponses struct {
	ID           uuid.UUID `json:"id"`
	NameTh       string          `json:"name_th"`
	NameEn       string          `json:"name_en"`
	MinSpending  decimal.Decimal `json:"min_spending"`
	IsActive     bool            `json:"is_active"`
	DiscountRate decimal.Decimal `json:"discount_rate"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListTierServiceRequest) ([]*ListTierServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`tiers.svc.list.start`)

	data, page, err := s.db.ListTiers(ctx, &entitiesdto.ListTiersRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListTierServiceResponses
	for _, item := range data {
		temp := &ListTierServiceResponses{
			ID:           item.ID,
			NameTh:       item.NameTh,
			NameEn:       item.NameEn,
			MinSpending:  item.MinSpending,
			IsActive:     item.IsActive,
			DiscountRate: item.DiscountRate,
			CreatedAt:    item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`tiers.svc.list.copy`)
	return response, page, nil
}
