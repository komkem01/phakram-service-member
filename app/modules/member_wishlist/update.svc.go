package member_wishlist

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

type UpdateMemberWishlistService struct {
	MemberID        uuid.UUID        `json:"member_id"`
	ProductID       uuid.UUID        `json:"product_id"`
	Quantity        *int             `json:"quantity"`
	PricePerUnit    *decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount *decimal.Decimal `json:"total_item_amount"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateMemberWishlistService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_wishlist.svc.update.start`)

	data, err := s.db.GetMemberWishlistByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.MemberID != uuid.Nil {
		data.MemberID = req.MemberID
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
			ActionType:   "member_wishlist",
			ActionID:     &data.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated member wishlist " + data.ID.String(),
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

	span.AddEvent(`member_wishlist.svc.update.success`)
	return nil
}
