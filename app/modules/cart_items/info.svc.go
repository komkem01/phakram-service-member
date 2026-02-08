package cart_items

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InfoCartItemServiceResponses struct {
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

func (s *Service) InfoService(ctx context.Context, id uuid.UUID, memberID uuid.UUID, isAdmin bool) (*InfoCartItemServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`cart_items.svc.info.start`)

	var data *ent.CartItemEntity
	if isAdmin || memberID == uuid.Nil {
		item, err := s.db.GetCartItemByID(ctx, id)
		if err != nil {
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return nil, err
		}
		data = item
	} else {
		item := new(ent.CartItemEntity)
		err := s.bunDB.DB().NewSelect().
			Model(item).
			Join("JOIN carts ON carts.id = cart_items.cart_id").
			Where("cart_items.id = ?", id).
			Where("carts.member_id = ?", memberID).
			Scan(ctx)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, sql.ErrNoRows
			}
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return nil, err
		}
		data = item
	}

	resp := &InfoCartItemServiceResponses{
		ID:              data.ID,
		CartID:          data.CartID,
		ProductID:       data.ProductID,
		Quantity:        data.Quantity,
		PricePerUnit:    data.PricePerUnit,
		TotalItemAmount: data.TotalItemAmount,
		IsSelected:      data.IsSelected,
		CreatedAt:       data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`cart_items.svc.info.success`)
	return resp, nil
}
