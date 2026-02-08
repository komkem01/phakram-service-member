package categories

import (
	"context"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UpdateCategoryService struct {
	ParentID *uuid.UUID `json:"parent_id"`
	NameTh   string     `json:"name_th"`
	NameEn   string     `json:"name_en"`
	IsActive *bool      `json:"is_active"`
	MemberID uuid.UUID  `json:"member_id"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdateCategoryService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`categories.svc.update.start`)

	data, err := s.db.GetCategoryByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.ParentID != nil {
		data.ParentID = req.ParentID
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

	nameForLog := data.NameEn
	if nameForLog == "" {
		nameForLog = data.NameTh
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
			ActionType:   "category",
			ActionID:     &data.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated category " + nameForLog,
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

	span.AddEvent(`categories.svc.update.success`)
	return nil
}
