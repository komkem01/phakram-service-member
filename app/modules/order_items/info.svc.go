package order_items

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InfoOrderItemServiceResponses struct {
	ID              uuid.UUID       `json:"id"`
	OrderID         uuid.UUID       `json:"order_id"`
	ProductID       uuid.UUID       `json:"product_id"`
	Quantity        int             `json:"quantity"`
	PricePerUnit    decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount decimal.Decimal `json:"total_item_amount"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID, memberID uuid.UUID, isAdmin bool) (*InfoOrderItemServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`order_items.svc.info.start`)

	var data *ent.OrderItemEntity
	if isAdmin || memberID == uuid.Nil {
		item, err := s.db.GetOrderItemByID(ctx, id)
		if err != nil {
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return nil, err
		}
		data = item
	} else {
		item := new(ent.OrderItemEntity)
		err := s.bunDB.DB().NewSelect().
			Model(item).
			Join("JOIN orders ON orders.id = order_items.order_id").
			Where("order_items.id = ?", id).
			Where("orders.member_id = ?", memberID).
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

	resp := &InfoOrderItemServiceResponses{
		ID:              data.ID,
		OrderID:         data.OrderID,
		ProductID:       data.ProductID,
		Quantity:        data.Quantity,
		PricePerUnit:    data.PricePerUnit,
		TotalItemAmount: data.TotalItemAmount,
		CreatedAt:       data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`order_items.svc.info.success`)
	return resp, nil
}
