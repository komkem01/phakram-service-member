package cart_items

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListCartItemServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

type ListCartItemServiceResponses struct {
	ID              uuid.UUID       `json:"id"`
	CartID          uuid.UUID       `json:"cart_id"`
	ProductID       uuid.UUID       `json:"product_id"`
	Quantity        int             `json:"quantity"`
	PricePerUnit    decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount decimal.Decimal `json:"total_item_amount"`
	IsSelected      bool            `json:"is_selected"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListCartItemServiceRequest) ([]*ListCartItemServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`cart_items.svc.list.start`)

	data, page, err := s.db.ListCartItems(ctx, &entitiesdto.ListCartItemsRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListCartItemServiceResponses
	for _, item := range data {
		temp := &ListCartItemServiceResponses{
			ID:              item.ID,
			CartID:          item.CartID,
			ProductID:       item.ProductID,
			Quantity:        item.Quantity,
			PricePerUnit:    item.PricePerUnit,
			TotalItemAmount: item.TotalItemAmount,
			IsSelected:      item.IsSelected,
			CreatedAt:       item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:       item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`cart_items.svc.list.copy`)
	return response, page, nil
}
