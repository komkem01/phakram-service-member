package prefixes

import (
	"context"
	"fmt"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreatePrefixService struct {
	NameTh   string    `json:"name_th"`
	NameEn   string    `json:"name_en"`
	GenderID uuid.UUID `json:"gender_id"`
	IsActive bool      `json:"is_active"`
}

func (s *Service) CreatePrefixService(ctx context.Context, req *CreatePrefixService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`prefixes.svc.create.start`)

	id := uuid.New()

	prefix := &ent.PrefixEntity{
		ID:       id,
		NameTh:   req.NameTh,
		NameEn:   req.NameEn,
		GenderID: req.GenderID,
		IsActive: req.IsActive,
	}
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(prefix).Exec(ctx); err != nil {
			return err
		}
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_prefix",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Created prefix with ID " + id.String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		span.AddEvent(`prefixes.svc.create.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_prefix",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Create prefix failed: %v", err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}
	span.AddEvent(`prefixes.svc.create.prefix_created`)

	span.AddEvent(`prefixes.svc.create.success`)
	return nil
}
