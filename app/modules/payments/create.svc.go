package payments

import (
	"context"
	"fmt"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type CreatePaymentService struct {
	Amount string `json:"amount"`
	Status string `json:"status"`
}

func (s *Service) CreatePaymentService(ctx context.Context, req *CreatePaymentService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`payments.svc.create.start`)

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return err
	}

	id := uuid.New()
	status, err := parsePaymentStatus(req.Status)
	if err != nil {
		return err
	}

	payment := &ent.PaymentEntity{ID: id, Amount: amount, Status: status}
	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(payment).Exec(ctx); err != nil {
			return err
		}
		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_payment",
			ActionID:     id,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Created payment with ID " + id.String(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		failLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_payment",
			ActionID:     id,
			Status:       ent.StatusAuditFailed,
			ActionDetail: fmt.Sprintf("Create payment failed: %v", err),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, _ = s.bunDB.DB().NewInsert().Model(failLog).Exec(ctx)
		return err
	}

	span.AddEvent(`payments.svc.create.success`)
	return nil
}

func parsePaymentStatus(status string) (ent.PaymentTypeEnum, error) {
	switch status {
	case "", string(ent.PaymentTypePending):
		return ent.PaymentTypePending, nil
	case string(ent.PaymentTypeSuccess):
		return ent.PaymentTypeSuccess, nil
	case string(ent.PaymentTypeFailed):
		return ent.PaymentTypeFailed, nil
	case string(ent.PaymentTypeRefunded):
		return ent.PaymentTypeRefunded, nil
	default:
		return "", fmt.Errorf("invalid payment status")
	}
}
