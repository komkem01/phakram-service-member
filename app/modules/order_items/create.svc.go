package order_items

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type CreateOrderItemService struct {
	OrderID         uuid.UUID       `json:"order_id"`
	ProductID       uuid.UUID       `json:"product_id"`
	Quantity        int             `json:"quantity"`
	PricePerUnit    decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount decimal.Decimal `json:"total_item_amount"`
	MemberID        uuid.UUID       `json:"member_id"`
}

func (s *Service) CreateOrderItemService(ctx context.Context, req *CreateOrderItemService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`order_items.svc.create.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		orderItem := &ent.OrderItemEntity{
			ID:              uuid.New(),
			OrderID:         req.OrderID,
			ProductID:       req.ProductID,
			Quantity:        req.Quantity,
			PricePerUnit:    req.PricePerUnit,
			TotalItemAmount: req.TotalItemAmount,
		}
		if _, err := tx.NewInsert().Model(orderItem).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "order_item",
			ActionID:     &orderItem.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created order item " + orderItem.ID.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	span.AddEvent(`order_items.svc.create.success`)
	return nil
}
