package cart_items

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func (s *Service) DeleteService(ctx context.Context, id uuid.UUID, memberID uuid.UUID) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`cart_items.svc.delete.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		item := new(ent.CartItemEntity)
		if err := tx.NewSelect().
			Model(item).
			Where("id = ?", id).
			Limit(1).
			Scan(ctx); err != nil {
			return err
		}

		if _, err := tx.NewDelete().Model(&ent.CartItemEntity{}).Where("id = ?", id).Exec(ctx); err != nil {
			return err
		}

		actionBy := memberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditDelete,
			ActionType:   "cart_item",
			ActionID:     &id,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Deleted cart item " + id.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		count, err := tx.NewSelect().
			Model((*ent.CartItemEntity)(nil)).
			Where("cart_id = ?", item.CartID).
			Count(ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		if count == 0 {
			if _, err := tx.NewUpdate().
				Model((*ent.CartEntity)(nil)).
				Set("is_active = ?", false).
				Where("id = ?", item.CartID).
				Exec(ctx); err != nil {
				return err
			}

			actionBy := memberID
			now := time.Now()
			cartAudit := &ent.AuditLogEntity{
				ID:           uuid.New(),
				Action:       ent.ActionAuditUpdate,
				ActionType:   "cart",
				ActionID:     &item.CartID,
				ActionBy:     &actionBy,
				Status:       ent.StatusAuditSuccess,
				ActionDetail: "Deactivated cart " + item.CartID.String(),
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			if _, err := tx.NewInsert().Model(cartAudit).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	span.AddEvent(`cart_items.svc.delete.success`)
	return nil
}
