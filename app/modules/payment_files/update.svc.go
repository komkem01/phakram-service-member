package payment_files

import (
	"context"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UpdatePaymentFileService struct {
	PaymentID uuid.UUID `json:"payment_id"`
	FileID    uuid.UUID `json:"file_id"`
	MemberID  uuid.UUID `json:"member_id"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdatePaymentFileService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`payment_files.svc.update.start`)

	data, err := s.db.GetPaymentFileByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.PaymentID != uuid.Nil {
		data.PaymentID = req.PaymentID
	}
	if req.FileID != uuid.Nil {
		data.FileID = req.FileID
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
			ActionType:   "payment_file",
			ActionID:     &data.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated payment file " + data.ID.String(),
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

	span.AddEvent(`payment_files.svc.update.success`)
	return nil
}
