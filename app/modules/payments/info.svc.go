package payments

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InfoPaymentServiceResponses struct {
	ID         uuid.UUID       `json:"id"`
	Amount     decimal.Decimal `json:"amount"`
	Status     string          `json:"status"`
	ApprovedBy uuid.UUID       `json:"approved_by"`
	ApprovedAt string          `json:"approved_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID, memberID uuid.UUID, isAdmin bool) (*InfoPaymentServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`payments.svc.info.start`)

	var data *ent.PaymentEntity
	if isAdmin || memberID == uuid.Nil {
		payment, err := s.db.GetPaymentByID(ctx, id)
		if err != nil {
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return nil, err
		}
		data = payment
	} else {
		payment := new(ent.PaymentEntity)
		err := s.bunDB.DB().NewSelect().
			Model(payment).
			Join("JOIN orders ON orders.payment_id = payments.id").
			Where("payments.id = ?", id).
			Where("orders.member_id = ?", memberID).
			Scan(ctx)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, sql.ErrNoRows
			}
			log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
			return nil, err
		}
		data = payment
	}

	resp := &InfoPaymentServiceResponses{
		ID:         data.ID,
		Amount:     data.Amount,
		Status:     string(data.Status),
		ApprovedBy: data.ApprovedBy,
		ApprovedAt: data.ApprovedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`payments.svc.info.success`)
	return resp, nil
}
