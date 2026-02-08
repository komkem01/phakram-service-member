package product_stocks

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

type UpdateProductStockService struct {
	ProductID   uuid.UUID        `json:"product_id"`
	UnitPrice   *decimal.Decimal `json:"unit_price"`
	StockAmount *int             `json:"stock_amount"`
	Remaining   *int             `json:"remaining"`
	MemberID    uuid.UUID        `json:"member_id"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateProductStockService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`product_stocks.svc.update.start`)

	data, err := s.db.GetProductStockByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.ProductID != uuid.Nil {
		data.ProductID = req.ProductID
	}
	if req.UnitPrice != nil {
		data.UnitPrice = *req.UnitPrice
	}
	if req.StockAmount != nil {
		data.StockAmount = *req.StockAmount
	}
	if req.Remaining != nil {
		data.Remaining = *req.Remaining
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
			ActionType:   "product_stock",
			ActionID:     &data.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated product stock " + data.ID.String(),
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

	span.AddEvent(`product_stocks.svc.update.success`)
	return nil
}
