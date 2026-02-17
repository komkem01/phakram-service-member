package products

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListProductServiceRequest struct {
	base.RequestPaginate
}

type ListProductServiceResponses struct {
	ID         uuid.UUID       `json:"id"`
	CategoryID uuid.UUID       `json:"category_id"`
	NameTh     string          `json:"name_th"`
	NameEn     string          `json:"name_en"`
	ProductNo  string          `json:"product_no"`
	Price      decimal.Decimal `json:"price"`
	IsActive   bool            `json:"is_active"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListProductServiceRequest) ([]*ListProductServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`products.svc.list.start`)

	data, page, err := s.db.ListProducts(ctx, &entitiesdto.ListProductsRequest{
		RequestPaginate: req.RequestPaginate,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListProductServiceResponses
	for _, item := range data {
		temp := &ListProductServiceResponses{
			ID:         item.ID,
			CategoryID: item.CategoryID,
			NameTh:     item.NameTh,
			NameEn:     item.NameEn,
			ProductNo:  item.ProductNo,
			Price:      item.Price,
			IsActive:   item.IsActive,
			CreatedAt:  item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`products.svc.list.copy`)
	return response, page, nil
}
