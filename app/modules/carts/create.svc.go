package carts

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreateCartService struct {
	MemberID uuid.UUID `json:"member_id"`
	IsActive bool      `json:"is_active"`
}

func (s *Service) CreateCartService(ctx context.Context, req *CreateCartService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.create.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		cart := &ent.CartEntity{
			ID:       uuid.New(),
			MemberID: req.MemberID,
			IsActive: req.IsActive,
		}
		if _, err := tx.NewInsert().Model(cart).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "cart",
			ActionID:     &cart.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created cart " + cart.ID.String(),
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
	span.AddEvent(`carts.svc.create.success`)
	return nil
}
