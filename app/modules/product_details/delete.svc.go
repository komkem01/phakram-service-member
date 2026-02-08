package product_details

import (
	"context"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func (s *Service) DeleteService(ctx context.Context, id uuid.UUID, memberID uuid.UUID) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`product_details.svc.delete.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model(&ent.ProductDetailEntity{}).Where("id = ?", id).Exec(ctx); err != nil {
			return err
		}

		actionBy := memberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditDelete,
			ActionType:   "product_detail",
			ActionID:     &id,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Deleted product detail " + id.String(),
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

	span.AddEvent(`product_details.svc.delete.success`)
	return nil
}
