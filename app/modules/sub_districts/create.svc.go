package sub_districts

import (
	"context"
	"fmt"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreateSubDistrictService struct {
	DistrictID uuid.UUID `json:"district_id"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
}

func (s *Service) CreateSubDistrictService(ctx context.Context, req *CreateSubDistrictService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`sub_districts.svc.create.start`)

	id := uuid.New()
	subDistrict := &ent.SubDistrictEntity{
		ID:         id,
		DistrictID: req.DistrictID,
		Name:       req.Name,
		IsActive:   req.IsActive,
	}
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(subDistrict).Exec(ctx); err != nil {
			return err
		}
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_sub_district",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Created sub district with ID " + id.String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		span.AddEvent(`sub_districts.svc.create.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_sub_district",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Create sub district failed: %v", err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}
	span.AddEvent(`sub_districts.svc.create.success`)
	return nil
}
