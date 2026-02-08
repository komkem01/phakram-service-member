package cart_items

import (
	"context"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type UpdateCartItemService struct {
	CartID          uuid.UUID        `json:"cart_id"`
	ProductID       uuid.UUID        `json:"product_id"`
	Quantity        *int             `json:"quantity"`
	PricePerUnit    *decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount *decimal.Decimal `json:"total_item_amount"`
	MemberID        uuid.UUID        `json:"member_id"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateCartItemService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`cart_items.svc.update.start`)

	data, err := s.db.GetCartItemByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.CartID != uuid.Nil {
		data.CartID = req.CartID
	}
	if req.ProductID != uuid.Nil {
		data.ProductID = req.ProductID
	}
	if req.Quantity != nil {
		data.Quantity = *req.Quantity
	}
	if req.PricePerUnit != nil {
		data.PricePerUnit = *req.PricePerUnit
	}
	if req.TotalItemAmount != nil {
		data.TotalItemAmount = *req.TotalItemAmount
	}

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditUpdate,
			ActionType:   "cart_item",
			ActionID:     &data.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated cart item " + data.ID.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`cart_items.svc.update.success`)
	return nil
}
