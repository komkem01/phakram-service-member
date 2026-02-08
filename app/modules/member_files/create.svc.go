package member_files

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreateMemberFileService struct {
	MemberID uuid.UUID `json:"member_id"`
	FileID   uuid.UUID `json:"file_id"`
}

func (s *Service) CreateMemberFileService(ctx context.Context, req *CreateMemberFileService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`member_files.svc.create.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		file := &ent.MemberFileEntity{
			ID:       uuid.New(),
			MemberID: req.MemberID,
			FileID:   req.FileID,
		}
		if _, err := tx.NewInsert().Model(file).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "member_file",
			ActionID:     &file.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created member file " + file.ID.String(),
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
	span.AddEvent(`member_files.svc.create.success`)
	return nil
}
