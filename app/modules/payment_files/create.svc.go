package payment_files

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CreatePaymentFileService struct {
	PaymentID uuid.UUID `json:"payment_id"`
	FileID    uuid.UUID `json:"file_id"`
	MemberID  uuid.UUID `json:"member_id"`
}

func (s *Service) CreatePaymentFileService(ctx context.Context, req *CreatePaymentFileService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`payment_files.svc.create.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		paymentFile := &ent.PaymentFileEntity{
			ID:        uuid.New(),
			PaymentID: req.PaymentID,
			FileID:    req.FileID,
		}
		if _, err := tx.NewInsert().Model(paymentFile).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "payment_file",
			ActionID:     &paymentFile.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created payment file " + paymentFile.ID.String(),
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
	span.AddEvent(`payment_files.svc.create.success`)
	return nil
}
