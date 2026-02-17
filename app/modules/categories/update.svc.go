package categories

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

type UpdateCategoryService struct {
	ParentID *string `json:"parent_id"`
	NameTh   string  `json:"name_th"`
	NameEn   string  `json:"name_en"`
	IsActive *bool   `json:"is_active"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateCategoryService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`categories.svc.update.start`)

	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		data := new(ent.CategoryEntity)
		if err := tx.NewSelect().Model(data).Where("id = ?", id).Scan(ctx); err != nil {
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return err
		}

		if req.ParentID != nil {
			if *req.ParentID == "" {
				data.ParentID = nil
			} else {
				parentID, err := uuid.Parse(*req.ParentID)
				if err != nil {
					return err
				}
				data.ParentID = &parentID
			}
		}
		if req.NameTh != "" {
			data.NameTh = req.NameTh
		}
		if req.NameEn != "" {
			data.NameEn = req.NameEn
		}
		if req.IsActive != nil {
			data.IsActive = *req.IsActive
		}

		if _, err := tx.NewUpdate().Model(data).Where("id = ?", data.ID).Exec(ctx); err != nil {
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_category",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Updated category with ID " + id.String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		span.AddEvent(`categories.svc.update.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_category",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Update category failed: %v", err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}

	span.AddEvent(`categories.svc.update.success`)
	return nil
}
