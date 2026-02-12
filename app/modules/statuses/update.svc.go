package statuses

import (
	"context"
	"fmt"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UpdateStatusService struct {
	NameTh string `json:"name_th"`
	NameEn string `json:"name_en"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateStatusService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`statuses.svc.update.start`)

	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.StatusEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", id).Scan(ctx); err != nil {
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return err
		}

		if req.NameTh != "" {
			data.NameTh = req.NameTh
		}
		if req.NameEn != "" {
			data.NameEn = req.NameEn
		}

		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_status",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Updated status with ID " + id.String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		span.AddEvent(`statuses.svc.update.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_status",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Update status failed: %v", err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}

	span.AddEvent(`statuses.svc.update.success`)
	return nil
}
