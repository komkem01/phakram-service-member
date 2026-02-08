package product_details

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListProductDetailServiceRequest struct {
	base.RequestPaginate
}

type ListProductDetailServiceResponses struct {
	ID               uuid.UUID       `json:"id"`
	ProductID        uuid.UUID       `json:"product_id"`
	Description      string          `json:"description"`
	Material         string          `json:"material"`
	Dimensions       string          `json:"dimensions"`
	Weight           decimal.Decimal `json:"weight"`
	CareInstructions string          `json:"care_instructions"`
}

func (s *Service) ListService(ctx context.Context, req *ListProductDetailServiceRequest) ([]*ListProductDetailServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`product_details.svc.list.start`)

	data, page, err := s.db.ListProductDetails(ctx, &entitiesdto.ListProductDetailsRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListProductDetailServiceResponses
	for _, item := range data {
		temp := &ListProductDetailServiceResponses{
			ID:               item.ID,
			ProductID:        item.ProductID,
			Description:      item.Description,
			Material:         item.Material,
			Dimensions:       item.Dimensions,
			Weight:           item.Weight,
			CareInstructions: item.CareInstructions,
		}
		response = append(response, temp)
	}
	span.AddEvent(`product_details.svc.list.copy`)
	return response, page, nil
}
