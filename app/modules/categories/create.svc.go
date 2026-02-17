package categories

import (
	"context"
	"fmt"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreateCategoryService struct {
	ParentID *string `json:"parent_id"`
	NameTh   string  `json:"name_th"`
	NameEn   string  `json:"name_en"`
	IsActive *bool   `json:"is_active"`
}

func (s *Service) CreateCategoryService(ctx context.Context, req *CreateCategoryService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`categories.svc.create.start`)

	id := uuid.New()
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	var parentID *uuid.UUID
	if req.ParentID != nil && *req.ParentID != "" {
		parsedParentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return err
		}
		parentID = &parsedParentID
	}

	category := &ent.CategoryEntity{
		ID:       id,
		ParentID: parentID,
		NameTh:   req.NameTh,
		NameEn:   req.NameEn,
		IsActive: isActive,
	}
	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(category).Exec(ctx); err != nil {
			return err
		}
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_category",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Created category with ID " + id.String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		span.AddEvent(`categories.svc.create.failed`)
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_category",
			ActionID:     id,
			ActionBy:     nil,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Create category failed: %v", err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}
	span.AddEvent(`categories.svc.create.success`)
	return nil
}
