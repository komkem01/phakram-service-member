package member_wishlist

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InfoMemberWishlistServiceResponses struct {
	ID              uuid.UUID       `json:"id"`
	MemberID        uuid.UUID       `json:"member_id"`
	ProductID       uuid.UUID       `json:"product_id"`
	Quantity        int             `json:"quantity"`
	PricePerUnit    decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount decimal.Decimal `json:"total_item_amount"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID, memberID uuid.UUID, isAdmin bool) (*InfoMemberWishlistServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_wishlist.svc.info.start`)

	data, err := s.db.GetMemberWishlistByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}
	if !isAdmin && memberID != uuid.Nil && data.MemberID != memberID {
		return nil, sql.ErrNoRows
	}

	resp := &InfoMemberWishlistServiceResponses{
		ID:              data.ID,
		MemberID:        data.MemberID,
		ProductID:       data.ProductID,
		Quantity:        data.Quantity,
		PricePerUnit:    data.PricePerUnit,
		TotalItemAmount: data.TotalItemAmount,
		CreatedAt:       data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`member_wishlist.svc.info.success`)
	return resp, nil
}
