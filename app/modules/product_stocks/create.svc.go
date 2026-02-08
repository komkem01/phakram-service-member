package product_stocks

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type CreateProductStockService struct {
	ProductID   uuid.UUID       `json:"product_id"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	StockAmount int             `json:"stock_amount"`
	Remaining   int             `json:"remaining"`
	MemberID    uuid.UUID       `json:"member_id"`
}

func (s *Service) CreateProductStockService(ctx context.Context, req *CreateProductStockService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`product_stocks.svc.create.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		productStock := &ent.ProductStockEntity{
			ID:          uuid.New(),
			ProductID:   req.ProductID,
			UnitPrice:   req.UnitPrice,
			StockAmount: req.StockAmount,
			Remaining:   req.Remaining,
		}
		if _, err := tx.NewInsert().Model(productStock).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "product_stock",
			ActionID:     &productStock.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created product stock " + productStock.ID.String(),
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
	span.AddEvent(`product_stocks.svc.create.success`)
	return nil
}
