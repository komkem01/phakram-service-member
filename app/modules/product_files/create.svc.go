package product_files

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreateProductFileService struct {
	ProductID uuid.UUID `json:"product_id"`
	FileID    uuid.UUID `json:"file_id"`
	MemberID  uuid.UUID `json:"member_id"`
}

func (s *Service) CreateProductFileService(ctx context.Context, req *CreateProductFileService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`product_files.svc.create.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		productFile := &ent.ProductFileEntity{
			ID:        uuid.New(),
			ProductID: req.ProductID,
			FileID:    req.FileID,
		}
		if _, err := tx.NewInsert().Model(productFile).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "product_file",
			ActionID:     &productFile.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created product file " + productFile.ID.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if _, err := tx.NewInsert().Model(auditLog).Exec(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	span.AddEvent(`product_files.svc.create.success`)
	return nil
}
