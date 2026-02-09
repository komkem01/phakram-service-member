package genders

import (
	"context"
	"fmt"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreateGenderService struct {
	NameTh   string `json:"name_th"`
	NameEn   string `json:"name_en"`
	IsActive bool   `json:"is_active"`
}

func (s *Service) CreateGenderService(ctx context.Context, req *CreateGenderService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`genders.svc.create.start`)

	id := uuid.New()

	// Create gender
	gender := &ent.GenderEntity{
		ID:       id,
		NameTh:   req.NameTh,
		NameEn:   req.NameEn,
		IsActive: req.IsActive,
	}
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(gender).Exec(ctx); err != nil {
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_gender",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Created gender with ID " + id.String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		span.AddEvent(`genders.svc.create.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_gender",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Create gender failed: %v", err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}
	span.AddEvent(`genders.svc.create.gender_created`)

	span.AddEvent(`genders.svc.create.success`)
	return nil
}
