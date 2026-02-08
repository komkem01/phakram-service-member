package order_items

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListOrderItemServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

type ListOrderItemServiceResponses struct {
	ID              uuid.UUID       `json:"id"`
	OrderID         uuid.UUID       `json:"order_id"`
	ProductID       uuid.UUID       `json:"product_id"`
	Quantity        int             `json:"quantity"`
	PricePerUnit    decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount decimal.Decimal `json:"total_item_amount"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListOrderItemServiceRequest) ([]*ListOrderItemServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`order_items.svc.list.start`)

	data, page, err := s.db.ListOrderItems(ctx, &entitiesdto.ListOrderItemsRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListOrderItemServiceResponses
	for _, item := range data {
		temp := &ListOrderItemServiceResponses{
			ID:              item.ID,
			OrderID:         item.OrderID,
			ProductID:       item.ProductID,
			Quantity:        item.Quantity,
			PricePerUnit:    item.PricePerUnit,
			TotalItemAmount: item.TotalItemAmount,
			CreatedAt:       item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:       item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`order_items.svc.list.copy`)
	return response, page, nil
}
