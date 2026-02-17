package payments

import (
	"context"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UpdatePaymentService struct {
	Amount string `json:"amount"`
	Status string `json:"status"`
}

func (s *Service) UpdatePaymentService(ctx context.Context, id string, req *UpdatePaymentService) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`payments.svc.update.start`)

	paymentID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	data, err := s.db.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return err
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return err
	}
	status, err := parsePaymentStatus(req.Status)
	if err != nil {
		return err
	}

	data.Amount = amount
	data.Status = status

	if err := s.db.UpdatePayment(ctx, data); err != nil {
		return err
	}

	span.AddEvent(`payments.svc.update.success`)
	return nil
}
