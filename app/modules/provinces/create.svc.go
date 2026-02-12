package provinces

import (
	"context"
	"fmt"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreateProvinceService struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func (s *Service) CreateProvinceService(ctx context.Context, req *CreateProvinceService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`provinces.svc.create.start`)

	id := uuid.New()
	province := &ent.ProvinceEntity{
		ID:       id,
		Name:     req.Name,
		IsActive: req.IsActive,
	}
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(province).Exec(ctx); err != nil {
			return err
		}
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_province",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Created province with ID " + id.String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		span.AddEvent(`provinces.svc.create.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_province",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Create province failed: %v", err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}
	span.AddEvent(`provinces.svc.create.success`)
	return nil
}
