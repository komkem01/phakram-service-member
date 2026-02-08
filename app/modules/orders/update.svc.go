package orders

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

type UpdateOrderService struct {
	OrderNo        string           `json:"order_no"`
	MemberID       uuid.UUID        `json:"member_id"`
	PaymentID      uuid.UUID        `json:"payment_id"`
	AddressID      uuid.UUID        `json:"address_id"`
	Status         string           `json:"status"`
	TotalAmount    *decimal.Decimal `json:"total_amount"`
	DiscountAmount *decimal.Decimal `json:"discount_amount"`
	NetAmount      *decimal.Decimal `json:"net_amount"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateOrderService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.update.start`)

	data, err := s.db.GetOrderByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.OrderNo != "" {
		data.OrderNo = req.OrderNo
	}
	if req.MemberID != uuid.Nil {
		data.MemberID = req.MemberID
	}
	if req.PaymentID != uuid.Nil {
		data.PaymentID = req.PaymentID
	}
	if req.AddressID != uuid.Nil {
		data.AddressID = req.AddressID
	}
	previousStatus := data.Status
	if req.Status != "" {
		data.Status = ent.StatusTypeEnum(req.Status)
	}
	if req.TotalAmount != nil {
		data.TotalAmount = *req.TotalAmount
	}
	if req.DiscountAmount != nil {
		data.DiscountAmount = *req.DiscountAmount
	}
	if req.NetAmount != nil {
		data.NetAmount = *req.NetAmount
	}

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		if actionBy == uuid.Nil {
			actionBy = data.MemberID
		}
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditUpdate,
			ActionType:   "order",
			ActionID:     &data.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated order " + data.ID.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		if previousStatus != data.Status && (data.Status == ent.StatusTypeCancelled || data.Status == ent.StatusTypeReturned) {
			orderItems := make([]*ent.OrderItemEntity, 0)
			if err := tx.NewSelect().
				Model(&orderItems).
				Where("order_id = ?", data.ID).
				Scan(ctx); err != nil {
				return err
			}

			for _, item := range orderItems {
				stock := new(ent.ProductStockEntity)
				if err := tx.NewSelect().
					Model(stock).
					Where("product_id = ?", item.ProductID).
					Limit(1).
					Scan(ctx); err != nil {
					return err
				}
				stock.Remaining = stock.Remaining + item.Quantity
				if _, err := tx.NewUpdate().
					Model(stock).
					Set("remaining = ?", stock.Remaining).
					Where("id = ?", stock.ID).
					Exec(ctx); err != nil {
					return err
				}

				stockAudit := &ent.AuditLogEntity{
					ID:           uuid.New(),
					Action:       ent.ActionAuditUpdate,
					ActionType:   "product_stock",
					ActionID:     &stock.ID,
					ActionBy:     &actionBy,
					Status:       ent.StatusAuditSuccess,
					ActionDetail: "Updated product stock " + stock.ID.String(),
					CreatedAt:    now,
					UpdatedAt:    now,
				}
				if _, err := tx.NewInsert().Model(stockAudit).Exec(ctx); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`orders.svc.update.success`)
	return nil
}
