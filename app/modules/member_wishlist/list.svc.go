package member_wishlist

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListMemberWishlistServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

type ListMemberWishlistServiceResponses struct {
	ID              uuid.UUID       `json:"id"`
	MemberID        uuid.UUID       `json:"member_id"`
	ProductID       uuid.UUID       `json:"product_id"`
	Quantity        int             `json:"quantity"`
	PricePerUnit    decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount decimal.Decimal `json:"total_item_amount"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListMemberWishlistServiceRequest) ([]*ListMemberWishlistServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_wishlist.svc.list.start`)

	data, page, err := s.db.ListMemberWishlist(ctx, &entitiesdto.ListMemberWishlistRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}
	var response []*ListMemberWishlistServiceResponses
	for _, item := range data {
		temp := &ListMemberWishlistServiceResponses{
			ID:              item.ID,
			MemberID:        item.MemberID,
			ProductID:       item.ProductID,
			Quantity:        item.Quantity,
			PricePerUnit:    item.PricePerUnit,
			TotalItemAmount: item.TotalItemAmount,
			CreatedAt:       item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:       item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response = append(response, temp)
	}
	span.AddEvent(`member_wishlist.svc.list.copy`)
	return response, page, nil
}
