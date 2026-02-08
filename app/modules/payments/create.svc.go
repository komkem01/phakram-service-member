package payments

import (
	"context"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type CreatePaymentService struct {
	Amount     decimal.Decimal `json:"amount"`
	Status     string          `json:"status"`
	ApprovedBy uuid.UUID       `json:"approved_by"`
	ApprovedAt time.Time       `json:"approved_at"`
	MemberID   uuid.UUID       `json:"member_id"`
}

func (s *Service) CreatePaymentService(ctx context.Context, req *CreatePaymentService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`payments.svc.create.start`)

	if err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		payment := &ent.PaymentEntity{
			ID:         uuid.New(),
			Amount:     req.Amount,
			Status:     ent.PaymentTypeEnum(req.Status),
			ApprovedBy: req.ApprovedBy,
			ApprovedAt: req.ApprovedAt,
		}
		if _, err := tx.NewInsert().Model(payment).Exec(ctx); err != nil {
			return err
		}

		actionBy := req.MemberID
		now := time.Now()
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.ActionAuditCreate,
			ActionType:   "payment",
			ActionID:     &payment.ID,
			ActionBy:     &actionBy,
			Status:       ent.StatusAuditSuccess,
			ActionDetail: "Created payment " + payment.ID.String(),
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
	span.AddEvent(`payments.svc.create.success`)
	return nil
}
