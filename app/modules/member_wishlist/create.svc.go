package member_wishlist

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type CreateMemberWishlistService struct {
	MemberID        uuid.UUID       `json:"member_id"`
	ProductID       uuid.UUID       `json:"product_id"`
	Quantity        int             `json:"quantity"`
	PricePerUnit    decimal.Decimal `json:"price_per_unit"`
	TotalItemAmount decimal.Decimal `json:"total_item_amount"`
}

func (s *Service) CreateMemberWishlistService(ctx context.Context, req *CreateMemberWishlistService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_wishlist.svc.create.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		wishlist := &ent.MemberWishlistEntity{
			ID:              uuid.New(),
			MemberID:        req.MemberID,
			ProductID:       req.ProductID,
			Quantity:        req.Quantity,
			PricePerUnit:    req.PricePerUnit,
			TotalItemAmount: req.TotalItemAmount,
		}
		if _, err := tx.NewInsert().Model(wishlist).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "member_wishlist",
			ActionID:     &wishlist.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created member wishlist " + wishlist.ID.String(),
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
	span.AddEvent(`member_wishlist.svc.create.success`)
	return nil
}
