package orders

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InfoOrderServiceResponses struct {
	ID             uuid.UUID       `json:"id"`
	OrderNo        string          `json:"order_no"`
	MemberID       uuid.UUID       `json:"member_id"`
	PaymentID      uuid.UUID       `json:"payment_id"`
	AddressID      uuid.UUID       `json:"address_id"`
	Status         string          `json:"status"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	NetAmount      decimal.Decimal `json:"net_amount"`
	CreatedAt      string          `json:"created_at"`
	UpdatedAt      string          `json:"updated_at"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID, memberID uuid.UUID, isAdmin bool) (*InfoOrderServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`orders.svc.info.start`)

	data, err := s.db.GetOrderByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}
	if !isAdmin && memberID != uuid.Nil && data.MemberID != memberID {
		return nil, sql.ErrNoRows
	}

	resp := &InfoOrderServiceResponses{
		ID:             data.ID,
		OrderNo:        data.OrderNo,
		MemberID:       data.MemberID,
		PaymentID:      data.PaymentID,
		AddressID:      data.AddressID,
		Status:         string(data.Status),
		TotalAmount:    data.TotalAmount,
		DiscountAmount: data.DiscountAmount,
		NetAmount:      data.NetAmount,
		CreatedAt:      data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	span.AddEvent(`orders.svc.info.success`)
	return resp, nil
}
