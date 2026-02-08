package payments

import (
	"context"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type UpdatePaymentService struct {
	Amount     *decimal.Decimal `json:"amount"`
	Status     string           `json:"status"`
	ApprovedBy uuid.UUID        `json:"approved_by"`
	ApprovedAt *time.Time       `json:"approved_at"`
	MemberID   uuid.UUID        `json:"member_id"`
}

func (s *Service) UpdateService(ctx context.Context, id uuid.UUID, req *UpdatePaymentService) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`payments.svc.update.start`)

	data, err := s.db.GetPaymentByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return err
	}

	if req.Amount != nil {
		data.Amount = *req.Amount
	}
	if req.Status != "" {
		data.Status = ent.PaymentTypeEnum(req.Status)
	}
	if req.ApprovedBy != uuid.Nil {
		data.ApprovedBy = req.ApprovedBy
	}
	if req.ApprovedAt != nil {
		data.ApprovedAt = *req.ApprovedAt
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
			ActionType:   "payment",
			ActionID:     &data.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Updated payment " + data.ID.String(),
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

	span.AddEvent(`payments.svc.update.success`)
	return nil
}
