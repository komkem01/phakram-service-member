package cart_items

import (
	"context"
	"database/sql"
	"errors"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type CreateCartItemService struct {
	CartID          uuid.UUID       `json:"cart_id"`
	ProductID       uuid.UUID       `json:"product_id"`
	Quantity        int             `json:"quantity"`
	PricePerUnit    decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount decimal.Decimal `json:"total_item_amount"`
	MemberID        uuid.UUID       `json:"member_id"`
}

func (s *Service) CreateCartItemService(ctx context.Context, req *CreateCartItemService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`cart_items.svc.create.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		cart := new(ent.CartEntity)
		if err := tx.NewSelect().
			Model(cart).
			Where("id = ?", req.CartID).
			Limit(1).
			Scan(ctx); err != nil {
			return err
		}

		if !cart.IsActive {
			if _, err := tx.NewUpdate().
				Model(cart).
				Set("is_active = ?", true).
				Where("id = ?", cart.ID).
				Exec(ctx); err != nil {
				return err
			}

			actionBy := req.MemberID
			now := time.Now()
			auditLog := &ent.AuditLogEntity{
				ID:           uuid.New(),
				Action:       ent.ActionAuditUpdate,
				ActionType:   "cart",
				ActionID:     &cart.ID,
				ActionBy:     &actionBy,
				Status:       ent.StatusAuditSuccess,
				ActionDetail: "Activated cart " + cart.ID.String(),
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
				return err
			}
		}

		existing := new(ent.CartItemEntity)
		err := tx.NewSelect().
			Model(existing).
			Where("cart_id = ?", req.CartID).
			Where("product_id = ?", req.ProductID).
			Limit(1).
			Scan(ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		if errors.Is(err, sql.ErrNoRows) {
			cartItem := &ent.CartItemEntity{
				ID:              uuid.New(),
				CartID:          req.CartID,
				ProductID:       req.ProductID,
				Quantity:        req.Quantity,
				PricePerUnit:    req.PricePerUnit,
				TotalItemAmount: req.TotalItemAmount,
			}
			if _, err := tx.NewInsert().Model(cartItem).Exec(ctx); err != nil {
				return err
			}

			auditLog := &ent.AuditLogEntity{
				ID:           uuid.New(),
				Action:       ent.ActionAuditCreate,
				ActionType:   "cart_item",
				ActionID:     &cartItem.ID,
				ActionBy:     &actionBy,
				Status:       ent.StatusAuditSuccess,
				ActionDetail: "Created cart item " + cartItem.ID.String(),
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
				return err
			}
			return nil
		}

		existing.Quantity += req.Quantity
		existing.PricePerUnit = req.PricePerUnit
		existing.TotalItemAmount = req.PricePerUnit.Mul(decimal.NewFromInt(int64(existing.Quantity)))
		if _, err := tx.NewUpdate().Model(existing).Where("id = ?", existing.ID).Exec(ctx); err != nil {
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditUpdate,
			ActionType:   "cart_item",
			ActionID:     &existing.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated cart item " + existing.ID.String(),
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
	span.AddEvent(`cart_items.svc.create.success`)
	return nil
}
