package carts

import (
	"context"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UpdateCartService struct {
	MemberID uuid.UUID `json:"member_id"`
	IsActive *bool     `json:"is_active"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateCartService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`carts.svc.update.start`)

	data, err := s.db.GetCartByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.MemberID != uuid.Nil {
		data.MemberID = req.MemberID
	}
	if req.IsActive != nil {
		data.IsActive = *req.IsActive
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
			ActionType:   "cart",
			ActionID:     &data.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated cart " + data.ID.String(),
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

	span.AddEvent(`carts.svc.update.success`)
	return nil
}
