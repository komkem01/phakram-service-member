package payments

import (
	"context"
	"errors"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
)

const errPaymentInUse = "payment is in use"

func (s *Service) DeletePaymentService(ctx context.Context, id string) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`payments.svc.delete.start`)

	paymentID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if err := s.canDeletePayment(ctx, paymentID); err != nil {
		return err
	}

	if err := s.db.DeletePayment(ctx, paymentID); err != nil {
		return err
	}

	span.AddEvent(`payments.svc.delete.success`)
	return nil
}

func (s *Service) canDeletePayment(ctx context.Context, paymentID uuid.UUID) error {
	orderCount, err := s.bunDB.DB().NewSelect().Model((*ent.OrderEntity)(nil)).Where("payment_id = ?", paymentID).Count(ctx)
	if err != nil {
		return err
	}
	if orderCount > 0 {
		return errors.New(errPaymentInUse)
	}

	memberPaymentCount, err := s.bunDB.DB().NewSelect().Model((*ent.MemberPaymentEntity)(nil)).Where("payment_id = ?", paymentID).Count(ctx)
	if err != nil {
		return err
	}
	if memberPaymentCount > 0 {
		return errors.New(errPaymentInUse)
	}

	paymentFileCount, err := s.bunDB.DB().NewSelect().Model((*ent.PaymentFileEntity)(nil)).Where("payment_id = ?", paymentID).Count(ctx)
	if err != nil {
		return err
	}
	if paymentFileCount > 0 {
		return errors.New(errPaymentInUse)
	}

	return nil
}
