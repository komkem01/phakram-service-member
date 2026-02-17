package payments

import (
	"context"
	"log/slog"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ListPaymentServiceRequest struct {
	base.RequestPaginate
	Status string
}

type ListPaymentServiceResponses struct {
	ID         uuid.UUID       `json:"id"`
	Amount     decimal.Decimal `json:"amount"`
	Status     string          `json:"status"`
	ApprovedBy *uuid.UUID      `json:"approved_by"`
	ApprovedAt *string         `json:"approved_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListPaymentServiceRequest) ([]*ListPaymentServiceResponses, *base.ResponsePaginate, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`payments.svc.list.start`)

	data, page, err := s.db.ListPayments(ctx, &entitiesdto.ListPaymentsRequest{RequestPaginate: req.RequestPaginate, Status: req.Status})
	if err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		return nil, nil, err
	}

	response := make([]*ListPaymentServiceResponses, 0, len(data))
	for _, item := range data {
		var approvedAt *string
		if item.ApprovedAt != nil {
			t := item.ApprovedAt.Format("2006-01-02T15:04:05Z07:00")
			approvedAt = &t
		}
		response = append(response, &ListPaymentServiceResponses{
			ID:         item.ID,
			Amount:     item.Amount,
			Status:     string(item.Status),
			ApprovedBy: item.ApprovedBy,
			ApprovedAt: approvedAt,
		})
	}

	span.AddEvent(`payments.svc.list.copy`)
	return response, page, nil
}
