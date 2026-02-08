package categories

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreateCategoryService struct {
	ParentID *uuid.UUID `json:"parent_id"`
	NameTh   string     `json:"name_th"`
	NameEn   string     `json:"name_en"`
	IsActive bool       `json:"is_active"`
	MemberID uuid.UUID  `json:"member_id"`
}

func (s *Service) CreateCategoryService(ctx context.Context, req *CreateCategoryService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`categories.svc.create.start`)

	// category := &ent.CategoryEntity{
	// 	ID:       uuid.New(),
	// 	ParentID: req.ParentID,
	// 	NameTh:   req.NameTh,
	// 	NameEn:   req.NameEn,
	// 	IsActive: req.IsActive,
	// }
	// if err := s.db.CreateCategory(ctx, category); err != nil {
	// 	return err
	// }

	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		category := &ent.CategoryEntity{
			ID:       uuid.New(),
			ParentID: req.ParentID,
			NameTh:   req.NameTh,
			NameEn:   req.NameEn,
			IsActive: req.IsActive,
		}
		if _, err := tx.NewInsert().Model(category).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "category",
			ActionID:     &category.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created category " + req.NameEn,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	span.AddEvent(`categories.svc.create.success`)
	return nil
}
